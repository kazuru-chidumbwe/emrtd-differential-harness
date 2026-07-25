package org.emrtd.harness.jmrtd;

/** Shared Observability Score contract (must match classifier/observability.go). */
public final class Observability {
    private Observability() {}

    public static final class TCAC01Outcome {
        public final boolean paceFailed;
        public final boolean bacSuccess;
        public final String bacErr;
        public final boolean paceSurfacedToCaller;

        public TCAC01Outcome(boolean paceFailed, boolean bacSuccess, String bacErr, boolean paceSurfacedToCaller) {
            this.paceFailed = paceFailed;
            this.bacSuccess = bacSuccess;
            this.bacErr = bacErr == null ? "" : bacErr;
            this.paceSurfacedToCaller = paceSurfacedToCaller;
        }
    }

    public static final class TCCA01Outcome {
        public final boolean chipAuthFailed;
        public final boolean chipAuthSuccess;
        public final boolean sessionContinueOk;
        public final boolean failureSurfacedToCaller;

        public TCCA01Outcome(boolean chipAuthFailed, boolean chipAuthSuccess,
                             boolean sessionContinueOk, boolean failureSurfacedToCaller) {
            this.chipAuthFailed = chipAuthFailed;
            this.chipAuthSuccess = chipAuthSuccess;
            this.sessionContinueOk = sessionContinueOk;
            this.failureSurfacedToCaller = failureSurfacedToCaller;
        }
    }

    public static final class TCAA01Outcome {
        public final boolean activeAuthFailed;
        public final boolean activeAuthSuccess;
        public final boolean failureSurfacedToCaller;

        public TCAA01Outcome(boolean activeAuthFailed, boolean activeAuthSuccess, boolean failureSurfacedToCaller) {
            this.activeAuthFailed = activeAuthFailed;
            this.activeAuthSuccess = activeAuthSuccess;
            this.failureSurfacedToCaller = failureSurfacedToCaller;
        }
    }

    public static final class TCTA01Outcome {
        public final boolean terminalAuthFailed;
        public final boolean terminalAuthSuccess;
        public final boolean failureSurfacedToCaller;
        public final boolean peerUnsupported;

        public TCTA01Outcome(boolean terminalAuthFailed, boolean terminalAuthSuccess,
                             boolean failureSurfacedToCaller, boolean peerUnsupported) {
            this.terminalAuthFailed = terminalAuthFailed;
            this.terminalAuthSuccess = terminalAuthSuccess;
            this.failureSurfacedToCaller = failureSurfacedToCaller;
            this.peerUnsupported = peerUnsupported;
        }
    }

    public static final class TCEAC01Outcome {
        public final boolean eacFailed;
        public final boolean eacSuccess;
        public final boolean protectedDgAccessible;
        public final boolean failureSurfacedToCaller;
        public final boolean peerUnsupported;

        public TCEAC01Outcome(boolean eacFailed, boolean eacSuccess, boolean protectedDgAccessible,
                              boolean failureSurfacedToCaller, boolean peerUnsupported) {
            this.eacFailed = eacFailed;
            this.eacSuccess = eacSuccess;
            this.protectedDgAccessible = protectedDgAccessible;
            this.failureSurfacedToCaller = failureSurfacedToCaller;
            this.peerUnsupported = peerUnsupported;
        }
    }

    public static int classifyTcAc01(TCAC01Outcome o) {
        if (o.paceFailed && o.bacSuccess && o.bacErr.isEmpty() && !o.paceSurfacedToCaller) {
            return 0;
        }
        if (o.paceFailed && !o.paceSurfacedToCaller) {
            return 1;
        }
        return 2;
    }

    public static int classifyTcCa01(TCCA01Outcome o) {
        if (o.chipAuthFailed && !o.chipAuthSuccess && o.sessionContinueOk && !o.failureSurfacedToCaller) {
            return 0;
        }
        if (o.chipAuthFailed && !o.failureSurfacedToCaller) {
            return 1;
        }
        return 2;
    }

    public static int classifyTcAa01(TCAA01Outcome o) {
        if (o.activeAuthFailed && !o.activeAuthSuccess && !o.failureSurfacedToCaller) {
            return 0;
        }
        if (o.activeAuthFailed && !o.failureSurfacedToCaller) {
            return 1;
        }
        return 2;
    }

    public static int classifyTcTa01(TCTA01Outcome o) {
        if (o.peerUnsupported) {
            return 2;
        }
        if (o.terminalAuthFailed && !o.terminalAuthSuccess && !o.failureSurfacedToCaller) {
            return 0;
        }
        if (o.terminalAuthFailed && !o.failureSurfacedToCaller) {
            return 1;
        }
        return 2;
    }

    public static int classifyTcEac01(TCEAC01Outcome o) {
        if (o.peerUnsupported) {
            return 2;
        }
        if (o.eacFailed && !o.eacSuccess && o.protectedDgAccessible && !o.failureSurfacedToCaller) {
            return 0;
        }
        if (o.eacFailed && !o.failureSurfacedToCaller) {
            return 1;
        }
        return 2;
    }

    public static String meaning(int score) {
        return switch (score) {
            case 0 -> "silent — failure not surfaced to caller";
            case 1 -> "logged — failure visible in session/trace only";
            default -> "surfaced — explicit error at caller boundary";
        };
    }
}
