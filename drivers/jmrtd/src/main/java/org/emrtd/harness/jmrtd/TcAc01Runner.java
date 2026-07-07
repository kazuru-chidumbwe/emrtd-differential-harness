package org.emrtd.harness.jmrtd;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import net.sf.scuba.smartcards.CardServiceException;
import net.sf.scuba.util.Hex;
import org.jmrtd.BACKey;
import org.jmrtd.PACEException;
import org.jmrtd.PassportService;
import org.jmrtd.lds.CardAccessFile;
import org.jmrtd.lds.PACEInfo;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.security.spec.AlgorithmParameterSpec;
import java.time.Instant;
import java.time.ZoneOffset;
import java.time.format.DateTimeFormatter;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

public final class TcAc01Runner {
    public static void main(String[] args) throws Exception {
        Path profilePath = Path.of("profiles/pace-then-bac-downgrade.json");
        Path logDir = Path.of("logs");
        for (int i = 0; i < args.length; i++) {
            if ("-profile".equals(args[i]) && i + 1 < args.length) {
                profilePath = Path.of(args[++i]);
            } else if ("-log-dir".equals(args[i]) && i + 1 < args.length) {
                logDir = Path.of(args[++i]);
            }
        }

        HarnessProfile profile = HarnessProfile.load(profilePath);
        String paceSw = profile.injection != null && profile.injection.paceSw != null
                ? profile.injection.paceSw : "6FFF";

        BACKey bacKey = new BACKey(
                profile.mrz.documentNumber,
                profile.mrz.dateOfBirth,
                profile.mrz.dateOfExpiry);

        TcAc01CardService card = new TcAc01CardService(bacKey, paceSw);
        PassportService service = new PassportService(card);
        service.open();

        CardAccessFile cardAccess = new CardAccessFile(
                new ByteArrayInputStream(Hex.hexStringToBytes(profile.cardAccessHex)));
        PACEInfo paceInfo = cardAccess.getPACEInfos().iterator().next();
        String oid = paceInfo.getObjectIdentifier();
        AlgorithmParameterSpec params = PACEInfo.toParameterSpec(paceInfo.getParameterId());

        String paceErr = "";
        boolean paceThrown = false;
        try {
            service.doPACE(bacKey, oid, params);
        } catch (PACEException e) {
            paceThrown = true;
            paceErr = e.getMessage();
        }

        service.sendSelectApplet(false);

        boolean bacSuccess = false;
        String bacErr = "";
        try {
            service.doBAC(bacKey);
            bacSuccess = true;
        } catch (CardServiceException e) {
            bacErr = e.getMessage();
        }

        int obs;
        String obsMeaning;
        if (paceThrown && bacSuccess && bacErr.isEmpty()) {
            // Naive integrator: catch PACEException and continue without rethrowing.
            obs = 0;
            obsMeaning = "silent — PACEException caught by harness integrator pattern; BAC succeeded without surfacing PACE failure upstream";
        } else if (paceThrown) {
            obs = 2;
            obsMeaning = "surfaced — PACEException thrown at API boundary";
        } else if (!paceErr.isEmpty()) {
            obs = 1;
            obsMeaning = "logged — PACE error recorded";
        } else {
            obs = 2;
            obsMeaning = "surfaced — explicit outcome";
        }

        String runId = profile.id + "-jmrtd-" + DateTimeFormatter.ofPattern("yyyyMMdd'T'HHmmss'Z'")
                .withZone(ZoneOffset.UTC).format(Instant.now());

        Map<String, Object> result = new LinkedHashMap<>();
        result.put("run_id", runId);
        result.put("test_case", profile.id);
        result.put("library", "jmrtd");
        result.put("mechanism", profile.mechanism);
        result.put("condition", profile.condition);
        result.put("pace_err", paceErr);
        result.put("pace_exception_thrown", paceThrown);
        result.put("bac_err", bacErr);
        result.put("bac_success", bacSuccess);
        result.put("observability_score", obs);
        result.put("observability_meaning", obsMeaning);

        List<Map<String, Object>> trace = new ArrayList<>();
        for (TcAc01CardService.TraceEntry e : card.trace()) {
            Map<String, Object> row = new LinkedHashMap<>();
            row.put("label", e.label);
            row.put("capdu", e.capdu);
            row.put("rapdu", e.rapdu);
            row.put("success", e.success);
            trace.add(row);
        }
        result.put("trace", trace);

        Files.createDirectories(logDir);
        Gson gson = new GsonBuilder().setPrettyPrinting().create();
        String json = gson.toJson(result);
        Path out = logDir.resolve(runId + ".json");
        Files.writeString(out, json);
        System.out.println(json);

        if (!paceThrown || !bacSuccess) {
            System.err.printf("TC-AC-01 gate failed: pace_thrown=%s bac_success=%s%n", paceThrown, bacSuccess);
            System.exit(1);
        }
    }
}
