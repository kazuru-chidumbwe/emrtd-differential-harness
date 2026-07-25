package org.emrtd.harness.jmrtd;

import net.sf.scuba.smartcards.CardServiceException;
import net.sf.scuba.util.Hex;
import org.jmrtd.BACKey;
import org.jmrtd.PACEException;
import org.jmrtd.PassportService;
import org.jmrtd.lds.CardAccessFile;
import org.jmrtd.lds.PACEInfo;

import java.io.ByteArrayInputStream;
import java.nio.file.Path;
import java.security.spec.AlgorithmParameterSpec;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

public final class TcAc01Runner {
    public static void main(String[] args) throws Exception {
        RunnerArgs a = RunnerArgs.parse(args);
        Path root = Path.of(".").toAbsolutePath().normalize();
        HarnessProfile profile = HarnessProfile.load(a.profilePath);
        String paceSw = profile.injection != null && profile.injection.paceSw != null
                ? profile.injection.paceSw : "6FFF";

        BACKey bacKey = new BACKey(profile.mrz.documentNumber, profile.mrz.dateOfBirth, profile.mrz.dateOfExpiry);
        TcAc01CardService card = new TcAc01CardService(bacKey, paceSw);
        PassportService service = PassportServices.open(card);
        service.open();

        CardAccessFile cardAccess = new CardAccessFile(new ByteArrayInputStream(Hex.hexStringToBytes(profile.cardAccessHex)));
        PACEInfo paceInfo = PassportServices.firstPaceInfo(cardAccess);
        String oid = paceInfo.getObjectIdentifier();
        AlgorithmParameterSpec params = PACEInfo.toParameterSpec(paceInfo.getParameterId());

        String paceErr = "";
        boolean paceThrown = false;
        try {
            service.doPACE(bacKey, oid, params, paceInfo.getParameterId());
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

        int obs = Observability.classifyTcAc01(new Observability.TCAC01Outcome(
                paceThrown || !paceErr.isEmpty(), bacSuccess, bacErr, false));

        Map<String, Object> result = new LinkedHashMap<>();
        result.put("run_id", RunIds.next(profile.id + "-jmrtd"));
        result.put("test_case", profile.id);
        result.put("library", "jmrtd");
        result.put("mechanism", profile.mechanism);
        result.put("condition", profile.condition);
        result.put("tier", profile.tier != null ? profile.tier : "wire");
        result.put("variant", a.variant);
        result.put("figure_id", a.figureId.isEmpty() ? null : a.figureId);
        result.put("pace_err", paceErr);
        result.put("pace_exception_thrown", paceThrown);
        result.put("bac_err", bacErr);
        result.put("bac_success", bacSuccess);
        result.put("observability_score", obs);
        result.put("observability_meaning", Observability.meaning(obs));
        result.put("provenance", Provenance.collect(root, a.profilePath, a.suiteId, a.suiteSeed, a.suiteN, a.runIndex,
                "java/TcAc01Runner", a.variant, ""));

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

        Provenance.writeResult(a.logDir, (String) result.get("run_id"), result);

        if (!paceThrown || !bacSuccess) {
            System.err.printf("TC-AC-01 gate failed: pace_thrown=%s bac_success=%s%n", paceThrown, bacSuccess);
            System.exit(1);
        }
    }
}
