package org.emrtd.harness.jmrtd;

import net.sf.scuba.smartcards.CardServiceException;
import org.jmrtd.BACKey;
import org.jmrtd.PassportService;

import java.nio.file.Path;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.interfaces.RSAPublicKey;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * TC-AA-01 baseline: BAC then doAA; chip rejects INTERNAL AUTHENTICATE.
 * Models a naive host that catches the exception and continues (FailureSurfacedToCaller=false).
 */
public final class TcAa01Runner {
    public static void main(String[] args) throws Exception {
        RunnerArgs a = RunnerArgs.parse(args);
        Path root = Path.of(".").toAbsolutePath().normalize();
        HarnessProfile profile = HarnessProfile.load(a.profilePath);
        String aaSw = profile.aaInjection != null && profile.aaInjection.aaSw != null
                ? profile.aaInjection.aaSw : "6982";

        BACKey bacKey = new BACKey(profile.mrz.documentNumber, profile.mrz.dateOfBirth, profile.mrz.dateOfExpiry);
        TcAa01CardService card = new TcAa01CardService(bacKey, aaSw);
        PassportService service = PassportServices.open(card);
        service.open();

        boolean bacSuccess = false;
        String bacErr = "";
        try {
            service.doBAC(bacKey);
            bacSuccess = true;
        } catch (CardServiceException e) {
            bacErr = e.getMessage();
        }

        String aaErr = "";
        boolean aaSuccess = false;
        if (bacSuccess) {
            KeyPairGenerator kpg = KeyPairGenerator.getInstance("RSA");
            kpg.initialize(2048);
            KeyPair kp = kpg.generateKeyPair();
            RSAPublicKey pub = (RSAPublicKey) kp.getPublic();
            byte[] challenge = new byte[] {1, 2, 3, 4, 5, 6, 7, 8};
            try {
                service.doAA(pub, "SHA-256", "SHA256withRSA", challenge);
                aaSuccess = true;
            } catch (Exception e) {
                // Naive host: catch and record; do not surface as hard stop.
                aaErr = e.getMessage() == null ? e.getClass().getSimpleName() : e.getMessage();
            }
        }

        int obs = Observability.classifyTcAa01(new Observability.TCAA01Outcome(
                !aaErr.isEmpty() || !aaSuccess, aaSuccess, false));

        Map<String, Object> result = new LinkedHashMap<>();
        result.put("run_id", RunIds.next(profile.id + "-jmrtd"));
        result.put("test_case", profile.id);
        result.put("library", "jmrtd");
        result.put("mechanism", profile.mechanism);
        result.put("condition", profile.condition);
        result.put("tier", profile.tier != null ? profile.tier : "wire");
        result.put("variant", a.variant);
        result.put("figure_id", a.figureId.isEmpty() ? null : a.figureId);
        result.put("active_auth_err", aaErr);
        result.put("active_auth_success", aaSuccess);
        result.put("bac_success", bacSuccess);
        result.put("bac_err", bacErr);
        result.put("observability_score", obs);
        result.put("observability_meaning", Observability.meaning(obs));
        result.put("provenance", Provenance.collect(root, a.profilePath, a.suiteId, a.suiteSeed, a.suiteN, a.runIndex,
                "java/TcAa01Runner", a.variant, ""));

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

        if (!bacSuccess || aaSuccess || aaErr.isEmpty()) {
            System.err.printf("TC-AA-01 gate failed: bac_success=%s aa_success=%s aa_err_empty=%s%n",
                    bacSuccess, aaSuccess, aaErr.isEmpty());
            System.exit(1);
        }
    }
}
