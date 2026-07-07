package org.emrtd.harness.jmrtd;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.security.MessageDigest;
import java.util.LinkedHashMap;
import java.util.Map;

public final class Provenance {
    private Provenance() {}

    public static Map<String, Object> collect(
            Path root,
            Path profilePath,
            String suiteId,
            int suiteSeed,
            int suiteN,
            int runIndex,
            String driver,
            String variant,
            String middleware) throws Exception {
        Map<String, Object> out = new LinkedHashMap<>();
        out.put("harness_commit", gitHead(root));
        out.put("harness_dirty", gitDirty(root));
        out.put("profile_path", profilePath.toString().replace('\\', '/'));
        out.put("profile_sha256", sha256(profilePath));
        out.put("suite_id", suiteId.isEmpty() ? null : suiteId);
        out.put("suite_seed", suiteSeed);
        out.put("suite_n", suiteN);
        out.put("run_index", runIndex);
        out.put("driver", driver);
        out.put("variant", variant);
        out.put("middleware", middleware == null || middleware.isEmpty() ? null : middleware);
        out.put("captured_at_utc", java.time.Instant.now().toString());
        return out;
    }

    private static String gitHead(Path root) {
        try {
            Process p = new ProcessBuilder("git", "-C", root.toString(), "rev-parse", "HEAD").start();
            return new String(p.getInputStream().readAllBytes()).trim();
        } catch (IOException | InterruptedException e) {
            Thread.currentThread().interrupt();
            return "unknown";
        }
    }

    private static boolean gitDirty(Path root) {
        try {
            Process p = new ProcessBuilder("git", "-C", root.toString(), "status", "--porcelain").start();
            return !new String(p.getInputStream().readAllBytes()).trim().isEmpty();
        } catch (IOException | InterruptedException e) {
            Thread.currentThread().interrupt();
            return false;
        }
    }

    private static String sha256(Path path) throws Exception {
        MessageDigest md = MessageDigest.getInstance("SHA-256");
        md.update(Files.readAllBytes(path));
        StringBuilder sb = new StringBuilder();
        for (byte b : md.digest()) {
            sb.append(String.format("%02x", b));
        }
        return sb.toString();
    }

    public static void writeResult(Path logDir, String runId, Map<String, Object> result) throws IOException {
        Files.createDirectories(logDir);
        Gson gson = new GsonBuilder().setPrettyPrinting().create();
        Files.writeString(logDir.resolve(runId + ".json"), gson.toJson(result));
        System.out.println(gson.toJson(result));
    }
}
