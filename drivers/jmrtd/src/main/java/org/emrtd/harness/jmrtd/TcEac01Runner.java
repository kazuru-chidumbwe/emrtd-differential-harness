package org.emrtd.harness.jmrtd;

import net.sf.scuba.smartcards.CardServiceException;
import net.sf.scuba.smartcards.CommandAPDU;
import net.sf.scuba.smartcards.ResponseAPDU;
import org.jmrtd.BACKey;
import org.jmrtd.PassportService;

import java.nio.file.Path;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * TC-EAC-01 baseline: TA fails but READ BINARY of protected DG still succeeds (silent downgrade).
 */
public final class TcEac01Runner {
    public static void main(String[] args) throws Exception {
        RunnerArgs a = RunnerArgs.parse(args);
        Path root = Path.of(".").toAbsolutePath().normalize();
        HarnessProfile profile = HarnessProfile.load(a.profilePath);
        String taSw = profile.taInjection != null && profile.taInjection.taSw != null
                ? profile.taInjection.taSw : "6982";

        BACKey bacKey = new BACKey(profile.mrz.documentNumber, profile.mrz.dateOfBirth, profile.mrz.dateOfExpiry);
        TcTa01CardService card = new TcTa01CardService(bacKey, taSw, true);
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

        String eacErr = "";
        boolean eacSuccess = false;
        boolean protectedDgAccessible = false;
        if (bacSuccess) {
            card.transmit(new CommandAPDU(0x00, 0x2A, 0x00, 0xBE, new byte[] {0x01, 0x02}, 0));
            boolean taOk = card.trace().stream()
                    .anyMatch(t -> t.label.equals("PSO:Verify Certificate (TA)") && t.success);
            if (!taOk) {
                eacErr = "EAC/TA PSO:Verify Certificate failed (synthetic)";
            } else {
                eacSuccess = true;
            }
            ResponseAPDU dg = card.transmit(new CommandAPDU(0x00, 0xB0, 0x00, 0x00, 8));
            protectedDgAccessible = dg.getSW() == 0x9000;
        }

        int obs = Observability.classifyTcEac01(new Observability.TCEAC01Outcome(
                !eacErr.isEmpty() || !eacSuccess,
                eacSuccess,
                protectedDgAccessible,
                false,
                false));

        Map<String, Object> result = new LinkedHashMap<>();
        result.put("run_id", RunIds.next(profile.id + "-jmrtd"));
        result.put("test_case", profile.id);
        result.put("library", "jmrtd");
        result.put("mechanism", profile.mechanism);
        result.put("condition", profile.condition);
        result.put("tier", profile.tier != null ? profile.tier : "wire");
        result.put("variant", a.variant);
        result.put("figure_id", a.figureId.isEmpty() ? null : a.figureId);
        result.put("eac_err", eacErr);
        result.put("eac_success", eacSuccess);
        result.put("protected_dg_accessible", protectedDgAccessible);
        result.put("bac_success", bacSuccess);
        result.put("bac_err", bacErr);
        result.put("peer_support", profile.peerSupport);
        result.put("observability_score", obs);
        result.put("observability_meaning", Observability.meaning(obs));
        result.put("provenance", Provenance.collect(root, a.profilePath, a.suiteId, a.suiteSeed, a.suiteN, a.runIndex,
                "java/TcEac01Runner", a.variant, ""));

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

        if (!bacSuccess || eacSuccess || !protectedDgAccessible || eacErr.isEmpty()) {
            System.err.printf("TC-EAC-01 gate failed: bac=%s eac_ok=%s dg=%s eac_err_empty=%s%n",
                    bacSuccess, eacSuccess, protectedDgAccessible, eacErr.isEmpty());
            System.exit(1);
        }
    }
}
