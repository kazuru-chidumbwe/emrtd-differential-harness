package org.emrtd.harness.jmrtd;

import net.sf.scuba.smartcards.CardServiceException;
import net.sf.scuba.smartcards.CommandAPDU;
import org.jmrtd.BACKey;
import org.jmrtd.PassportService;

import java.nio.file.Path;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * TC-CA-01 mitigated: explicit-reject on chip authentication failure (JMRTD analogue of
 * middleware/ca.go's PerformChipAuth with AllowContinue=false). A session is not treated as
 * usable if chip authentication was attempted and failed; the failure is a hard, surfaced
 * error rather than a value the caller can silently ignore alongside BAC success.
 */
public final class TcCa01MitigatedRunner {
    public static void main(String[] args) throws Exception {
        RunnerArgs a = RunnerArgs.parse(args);
        Path root = Path.of(".").toAbsolutePath().normalize();
        HarnessProfile profile = HarnessProfile.load(a.profilePath);
        String caSw = profile.caInjection != null && profile.caInjection.caSw != null
                ? profile.caInjection.caSw : "6FFF";

        BACKey bacKey = new BACKey(profile.mrz.documentNumber, profile.mrz.dateOfBirth, profile.mrz.dateOfExpiry);
        TcCa01CardService card = new TcCa01CardService(bacKey, caSw);
        PassportService service = new PassportService(card);
        service.open();

        boolean bacSuccess = false;
        String bacErr = "";
        try {
            service.doBAC(bacKey);
            bacSuccess = true;
        } catch (CardServiceException e) {
            bacErr = e.getMessage();
        }

        String chipAuthErr = "";
        boolean chipAuthSuccess = false;
        String middlewareErr = "";
        if (bacSuccess) {
            card.transmit(new CommandAPDU(0x00, 0x22, 0x41, 0xA4, new byte[0], 0));
            chipAuthSuccess = card.trace().stream()
                    .anyMatch(t -> t.label.equals("MSE:Set AT (CA)") && t.success);
            if (!chipAuthSuccess) {
                chipAuthErr = "CA MSE:Set AT failed (synthetic chip)";
                // Explicit-reject: chip authentication failure is a hard stop. The session is
                // not handed back as "open" the way the baseline runner implicitly allows.
                middlewareErr = "chip authentication failed: explicit reject (middleware CA analogue): " + chipAuthErr;
            }
        }

        boolean failureSurfaced = !middlewareErr.isEmpty();
        int obs = Observability.classifyTcCa01(new Observability.TCCA01Outcome(
                !chipAuthSuccess, chipAuthSuccess, failureSurfaced));

        Map<String, Object> result = new LinkedHashMap<>();
        result.put("run_id", RunIds.next(profile.id + "-jmrtd-mitigated"));
        result.put("test_case", profile.id);
        result.put("library", "jmrtd");
        result.put("mechanism", profile.mechanism);
        result.put("condition", profile.condition);
        result.put("tier", profile.tier != null ? profile.tier : "wire");
        result.put("variant", a.variant);
        result.put("figure_id", a.figureId.isEmpty() ? null : a.figureId);
        result.put("chip_auth_err", chipAuthErr);
        result.put("chip_auth_success", chipAuthSuccess);
        result.put("bac_success", bacSuccess);
        result.put("bac_err", bacErr);
        result.put("middleware_err", middlewareErr);
        result.put("observability_score", obs);
        result.put("observability_meaning", Observability.meaning(obs));
        result.put("provenance", Provenance.collect(root, a.profilePath, a.suiteId, a.suiteSeed, a.suiteN, a.runIndex,
                "java/TcCa01MitigatedRunner", a.variant, "explicit-reject-ca"));

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

        // Gate: the synthetic profile is expected to fail CA. A mitigated run that doesn't
        // surface that failure (middlewareErr empty) or that reports CA succeeding indicates
        // the harness fixture or the mitigation wiring is broken, not a real finding.
        if (bacSuccess && chipAuthSuccess) {
            System.err.println("TC-CA-01 mitigated gate failed: expected CA failure, got success");
            System.exit(1);
        }
        if (bacSuccess && middlewareErr.isEmpty()) {
            System.err.println("TC-CA-01 mitigated gate failed: CA failed but was not surfaced");
            System.exit(1);
        }
    }
}
