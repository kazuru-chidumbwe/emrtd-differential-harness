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
 * TC-TA-01 baseline (JMRTD-asymmetric): BAC then PSO:Verify Certificate fails.
 * Naive host records failure without surfacing (score 0).
 */
public final class TcTa01Runner {
    public static void main(String[] args) throws Exception {
        RunnerArgs a = RunnerArgs.parse(args);
        Path root = Path.of(".").toAbsolutePath().normalize();
        HarnessProfile profile = HarnessProfile.load(a.profilePath);
        String taSw = profile.taInjection != null && profile.taInjection.taSw != null
                ? profile.taInjection.taSw : "6982";

        BACKey bacKey = new BACKey(profile.mrz.documentNumber, profile.mrz.dateOfBirth, profile.mrz.dateOfExpiry);
        TcTa01CardService card = new TcTa01CardService(bacKey, taSw, false);
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

        String taErr = "";
        boolean taSuccess = false;
        if (bacSuccess) {
            // SW-proxy: send PSO:Verify Certificate (INS 0x2A, P1=0x00, P2=0xBE)
            card.transmit(new CommandAPDU(0x00, 0x2A, 0x00, 0xBE, new byte[] {0x01, 0x02}, 0));
            taSuccess = card.trace().stream()
                    .anyMatch(t -> t.label.equals("PSO:Verify Certificate (TA)") && t.success);
            if (!taSuccess) {
                taErr = "TA PSO:Verify Certificate failed (synthetic chip)";
            }
        }

        int obs = Observability.classifyTcTa01(new Observability.TCTA01Outcome(
                !taErr.isEmpty() || !taSuccess, taSuccess, false, false));

        Map<String, Object> result = new LinkedHashMap<>();
        result.put("run_id", RunIds.next(profile.id + "-jmrtd"));
        result.put("test_case", profile.id);
        result.put("library", "jmrtd");
        result.put("mechanism", profile.mechanism);
        result.put("condition", profile.condition);
        result.put("tier", profile.tier != null ? profile.tier : "wire");
        result.put("variant", a.variant);
        result.put("figure_id", a.figureId.isEmpty() ? null : a.figureId);
        result.put("terminal_auth_err", taErr);
        result.put("terminal_auth_success", taSuccess);
        result.put("bac_success", bacSuccess);
        result.put("bac_err", bacErr);
        result.put("peer_support", profile.peerSupport);
        result.put("observability_score", obs);
        result.put("observability_meaning", Observability.meaning(obs));
        result.put("provenance", Provenance.collect(root, a.profilePath, a.suiteId, a.suiteSeed, a.suiteN, a.runIndex,
                "java/TcTa01Runner", a.variant, ""));

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

        if (!bacSuccess || taSuccess || taErr.isEmpty()) {
            System.err.printf("TC-TA-01 gate failed: bac=%s ta_success=%s ta_err_empty=%s%n",
                    bacSuccess, taSuccess, taErr.isEmpty());
            System.exit(1);
        }
    }
}
