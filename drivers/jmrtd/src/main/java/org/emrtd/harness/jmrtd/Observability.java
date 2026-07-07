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

    public static int classifyTcAc01(TCAC01Outcome o) {
        if (o.paceFailed && o.bacSuccess && o.bacErr.isEmpty() && !o.paceSurfacedToCaller) {
            return 0;
        }
        if (o.paceFailed && !o.paceSurfacedToCaller) {
            return 1;
        }
        return 2;
    }

    public static String meaningTcAc01(int score) {
        return switch (score) {
            case 0 -> "silent — PACE failure not surfaced to caller; BAC succeeded";
            case 1 -> "logged — PACE failure visible in session/trace only";
            default -> "surfaced — explicit error at caller boundary";
        };
    }
}
