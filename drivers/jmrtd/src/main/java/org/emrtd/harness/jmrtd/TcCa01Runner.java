package org.emrtd.harness.jmrtd;

import net.sf.scuba.smartcards.CardServiceException;
import net.sf.scuba.smartcards.CommandAPDU;
import net.sf.scuba.smartcards.ResponseAPDU;
import org.jmrtd.BACKey;
import org.jmrtd.PassportService;
import org.jmrtd.protocol.SecureMessagingWrapper;

import java.nio.file.Path;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

public final class TcCa01Runner {
    public static void main(String[] args) throws Exception {
        RunnerArgs a = RunnerArgs.parse(args);
        Path root = Path.of(".").toAbsolutePath().normalize();
        HarnessProfile profile = HarnessProfile.load(a.profilePath);
        String caSw = profile.caInjection != null && profile.caInjection.caSw != null
                ? profile.caInjection.caSw : "6FFF";

        BACKey bacKey = new BACKey(profile.mrz.documentNumber, profile.mrz.dateOfBirth, profile.mrz.dateOfExpiry);
        TcCa01CardService card = new TcCa01CardService(bacKey, caSw);
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

        String chipAuthErr = "";
        boolean chipAuthSuccess = false;
        boolean sessionContinueOk = false;
        if (bacSuccess) {
            SecureMessagingWrapper wrapper = service.getWrapper();
            if (wrapper == null) {
                throw new IllegalStateException("PassportService.getWrapper() null after doBAC");
            }
            // CA MSE via BAC SM wrapper (not raw unprotected APDU).
            ResponseAPDU caRsp = wrapper.unwrap(card.transmit(wrapper.wrap(
                    new CommandAPDU(0x00, 0x22, 0x41, 0xA4, new byte[0], 0))));
            chipAuthSuccess = caRsp.getSW() == 0x9000;
            if (!chipAuthSuccess) {
                chipAuthErr = "CA MSE:Set AT failed (synthetic chip)";
                // Confirm BAC SM still usable after CA reject (bilateral SM continue-check).
                if (service.getWrapper() == null) {
                    throw new IllegalStateException("wrapper cleared after CA MSE reject");
                }
                ResponseAPDU dg = wrapper.unwrap(card.transmit(wrapper.wrap(
                        new CommandAPDU(0x00, 0xB0, 0x00, 0x00, 5))));
                sessionContinueOk = dg.getSW() == 0x9000 && dg.getData() != null && dg.getData().length > 0;
            }
        }

        int obs = Observability.classifyTcCa01(new Observability.TCCA01Outcome(
                !chipAuthErr.isEmpty() || !chipAuthSuccess, chipAuthSuccess, sessionContinueOk, false));

        Map<String, Object> result = new LinkedHashMap<>();
        result.put("run_id", RunIds.next(profile.id + "-jmrtd"));
        result.put("test_case", profile.id);
        result.put("library", "jmrtd");
        result.put("mechanism", profile.mechanism);
        result.put("condition", profile.condition);
        result.put("tier", profile.tier != null ? profile.tier : "wire");
        result.put("variant", a.variant);
        result.put("figure_id", a.figureId.isEmpty() ? null : a.figureId);
        result.put("chip_auth_err", chipAuthErr);
        result.put("chip_auth_success", chipAuthSuccess);
        result.put("session_continue_ok", sessionContinueOk);
        result.put("bac_success", bacSuccess);
        result.put("bac_err", bacErr);
        result.put("observability_score", obs);
        result.put("observability_meaning", Observability.meaning(obs));
        result.put("provenance", Provenance.collect(root, a.profilePath, a.suiteId, a.suiteSeed, a.suiteN, a.runIndex,
                "java/TcCa01Runner", a.variant, ""));

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

        if (!bacSuccess || chipAuthSuccess || !sessionContinueOk) {
            System.err.printf("TC-CA-01 gate failed: bac_success=%s chip_auth_success=%s session_continue=%s%n",
                    bacSuccess, chipAuthSuccess, sessionContinueOk);
            System.exit(1);
        }
    }
}
