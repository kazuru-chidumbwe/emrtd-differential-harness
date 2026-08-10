package org.emrtd.harness.jmrtd;

import org.junit.jupiter.api.Test;

import java.nio.file.Files;
import java.nio.file.Path;

import static org.junit.jupiter.api.Assertions.assertEquals;

/** Cross-language Observability Score contract tests (must match Go + Python). */
public class ObservabilityTest {

    @Test
    void classifyTcAc01Silent() {
        assertEquals(0, Observability.classifyTcAc01(
                new Observability.TCAC01Outcome(true, true, "", false)));
    }

    @Test
    void classifyTcAc01Logged() {
        assertEquals(1, Observability.classifyTcAc01(
                new Observability.TCAC01Outcome(true, false, "bac failed", false)));
    }

    @Test
    void classifyTcAc01Surfaced() {
        assertEquals(2, Observability.classifyTcAc01(
                new Observability.TCAC01Outcome(true, false, "", true)));
    }

    @Test
    void classifyTcCa01Silent() {
        assertEquals(0, Observability.classifyTcCa01(
                new Observability.TCCA01Outcome(true, false, true, false)));
    }

    @Test
    void classifyTcCa01Logged() {
        assertEquals(1, Observability.classifyTcCa01(
                new Observability.TCCA01Outcome(true, false, false, false)));
    }

    @Test
    void classifyTcCa01Surfaced() {
        assertEquals(2, Observability.classifyTcCa01(
                new Observability.TCCA01Outcome(true, false, true, true)));
    }

    @Test
    void classifyTcAa01Silent() {
        assertEquals(0, Observability.classifyTcAa01(
                new Observability.TCAA01Outcome(true, false, false)));
    }

    @Test
    void classifyTcAa01Logged() {
        assertEquals(1, Observability.classifyTcAa01(
                new Observability.TCAA01Outcome(true, true, false)));
    }

    @Test
    void classifyTcAa01Surfaced() {
        assertEquals(2, Observability.classifyTcAa01(
                new Observability.TCAA01Outcome(true, false, true)));
    }

    @Test
    void classifyTcTa01PeerUnsupported() {
        assertEquals(2, Observability.classifyTcTa01(
                new Observability.TCTA01Outcome(false, false, false, true)));
    }

    @Test
    void classifyTcTa01Silent() {
        assertEquals(0, Observability.classifyTcTa01(
                new Observability.TCTA01Outcome(true, false, false, false)));
    }

    @Test
    void classifyTcTa01Logged() {
        assertEquals(1, Observability.classifyTcTa01(
                new Observability.TCTA01Outcome(true, true, false, false)));
    }

    @Test
    void classifyTcTa01Surfaced() {
        assertEquals(2, Observability.classifyTcTa01(
                new Observability.TCTA01Outcome(true, false, true, false)));
    }

    @Test
    void classifyTcEac01PeerUnsupported() {
        assertEquals(2, Observability.classifyTcEac01(
                new Observability.TCEAC01Outcome(false, false, false, false, true)));
    }

    @Test
    void classifyTcEac01SilentProtectedDg() {
        assertEquals(0, Observability.classifyTcEac01(
                new Observability.TCEAC01Outcome(true, false, true, false, false)));
    }

    @Test
    void classifyTcEac01LoggedNoDg() {
        assertEquals(1, Observability.classifyTcEac01(
                new Observability.TCEAC01Outcome(true, false, false, false, false)));
    }

    @Test
    void classifyTcEac01Surfaced() {
        assertEquals(2, Observability.classifyTcEac01(
                new Observability.TCEAC01Outcome(true, false, true, true, false)));
    }

    @Test
    void sharedObservabilityVectors() throws Exception {
        Path vectors = Path.of("..", "..", "..", "..", "..", "testdata", "observability-vectors.json")
                .toAbsolutePath().normalize();
        if (!Files.isRegularFile(vectors)) {
            // Maven cwd is drivers/jmrtd
            vectors = Path.of("testdata", "observability-vectors.json");
            if (!Files.isRegularFile(vectors)) {
                vectors = Path.of("..", "..", "testdata", "observability-vectors.json").normalize();
            }
        }
        // Resolve from harness root relative to this class location
        Path here = Path.of("src/test/java").toAbsolutePath().normalize();
        Path harness = here;
        for (int i = 0; i < 8; i++) {
            if (Files.isRegularFile(harness.resolve("testdata/observability-vectors.json"))) {
                vectors = harness.resolve("testdata/observability-vectors.json");
                break;
            }
            harness = harness.getParent();
            if (harness == null) {
                break;
            }
        }
        String raw = Files.readString(vectors);
        com.google.gson.JsonArray arr = com.google.gson.JsonParser.parseString(raw).getAsJsonArray();
        for (com.google.gson.JsonElement el : arr) {
            com.google.gson.JsonObject row = el.getAsJsonObject();
            String id = row.get("id").getAsString();
            String mech = row.get("mechanism").getAsString();
            com.google.gson.JsonObject inp = row.getAsJsonObject("input");
            int want = row.get("expected_score").getAsInt();
            int got;
            switch (mech) {
                case "TC-AC-01" -> got = Observability.classifyTcAc01(new Observability.TCAC01Outcome(
                        inp.get("pace_failed").getAsBoolean(),
                        inp.get("bac_success").getAsBoolean(),
                        inp.has("bac_err") && !inp.get("bac_err").isJsonNull() ? inp.get("bac_err").getAsString() : "",
                        inp.get("pace_surfaced_to_caller").getAsBoolean()));
                case "TC-CA-01" -> got = Observability.classifyTcCa01(new Observability.TCCA01Outcome(
                        inp.get("chip_auth_failed").getAsBoolean(),
                        inp.get("chip_auth_success").getAsBoolean(),
                        inp.get("session_continue_ok").getAsBoolean(),
                        inp.get("failure_surfaced_to_caller").getAsBoolean()));
                case "TC-AA-01" -> got = Observability.classifyTcAa01(new Observability.TCAA01Outcome(
                        inp.get("active_auth_failed").getAsBoolean(),
                        inp.get("active_auth_success").getAsBoolean(),
                        inp.get("failure_surfaced_to_caller").getAsBoolean()));
                default -> throw new IllegalStateException("unsupported " + mech + " (" + id + ")");
            }
            assertEquals(want, got, id);
        }
    }
}
