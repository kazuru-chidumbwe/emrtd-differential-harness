package org.emrtd.harness.jmrtd;

import org.junit.jupiter.api.Test;

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
                new Observability.TCCA01Outcome(true, false, false)));
    }

    @Test
    void classifyTcCa01Logged() {
        assertEquals(1, Observability.classifyTcCa01(
                new Observability.TCCA01Outcome(true, true, false)));
    }

    @Test
    void classifyTcCa01Surfaced() {
        assertEquals(2, Observability.classifyTcCa01(
                new Observability.TCCA01Outcome(true, false, true)));
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
}
