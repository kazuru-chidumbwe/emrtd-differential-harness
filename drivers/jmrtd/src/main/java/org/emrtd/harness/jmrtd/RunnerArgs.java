package org.emrtd.harness.jmrtd;

import java.nio.file.Path;
import java.time.Instant;
import java.time.ZoneOffset;
import java.time.format.DateTimeFormatter;
import java.util.concurrent.atomic.AtomicLong;

final class RunIds {
    private static final AtomicLong SEQ = new AtomicLong();

    private RunIds() {}

    static String next(String prefix) {
        long n = SEQ.incrementAndGet();
        String ts = DateTimeFormatter.ofPattern("yyyyMMdd'T'HHmmss.SSSSSS'Z'")
                .withZone(ZoneOffset.UTC).format(Instant.now());
        return String.format("%s-%s-%06d", prefix, ts, n);
    }
}

final class RunnerArgs {
    Path profilePath = Path.of("profiles/pace-then-bac-downgrade.json");
    Path logDir = Path.of("logs");
    String variant = "baseline";
    String suiteId = "";
    int suiteSeed = 1;
    int suiteN = 1;
    int runIndex = 1;
    String figureId = "";

    static RunnerArgs parse(String[] args) {
        RunnerArgs a = new RunnerArgs();
        for (int i = 0; i < args.length; i++) {
            switch (args[i]) {
                case "-profile" -> a.profilePath = Path.of(args[++i]);
                case "-log-dir" -> a.logDir = Path.of(args[++i]);
                case "-variant" -> a.variant = args[++i];
                case "-suite-id" -> a.suiteId = args[++i];
                case "-suite-seed" -> a.suiteSeed = Integer.parseInt(args[++i]);
                case "-suite-n" -> a.suiteN = Integer.parseInt(args[++i]);
                case "-run-index" -> a.runIndex = Integer.parseInt(args[++i]);
                case "-figure-id" -> a.figureId = args[++i];
                default -> { }
            }
        }
        return a;
    }
}
