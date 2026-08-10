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
        out.put("harness_commit", resolveHarnessCommit(root));
        out.put("harness_dirty", gitDirty(root));
        out.put("profile_path", profilePath.toString().replace('\\', '/'));
        out.put("profile_sha256", sha256(profilePath));
        out.put("suite_id", suiteId == null || suiteId.isEmpty() ? "unspecified" : suiteId);
        out.put("suite_seed", suiteSeed);
        out.put("suite_n", suiteN);
        // 1-based suite index (run_suite.py uses range(1, n+1); never emit 0)
        out.put("run_index", runIndex < 1 ? 1 : runIndex);
        out.put("driver", driver);
        out.put("variant", variant);
        out.put("middleware", middleware == null || middleware.isEmpty() ? null : middleware);
        out.put("captured_at_utc", java.time.Instant.now().toString());
        return out;
    }

    /** Prefer EMRTD_HARNESS_COMMIT, then git rev-parse HEAD; never return empty. */
    static String resolveHarnessCommit(Path root) {
        String env = System.getenv("EMRTD_HARNESS_COMMIT");
        if (env != null) {
            String trimmed = env.trim();
            if (!trimmed.isEmpty()) {
                return trimmed;
            }
        }
        Path gitRoot = findGitRoot(root);
        if (gitRoot != null) {
            String head = gitRevParseHead(gitRoot);
            if (head != null && !head.isEmpty()) {
                return head;
            }
        }
        return "unknown";
    }

    private static Path findGitRoot(Path start) {
        Path dir = start.toAbsolutePath().normalize();
        for (int i = 0; i < 32; i++) {
            if (Files.isDirectory(dir.resolve(".git")) || Files.isRegularFile(dir.resolve(".git"))) {
                return dir;
            }
            Path parent = dir.getParent();
            if (parent == null || parent.equals(dir)) {
                return null;
            }
            dir = parent;
        }
        return null;
    }

    private static String gitRevParseHead(Path root) {
        try {
            Process p = new ProcessBuilder("git", "-C", root.toString(), "rev-parse", "HEAD")
                    .redirectErrorStream(true)
                    .start();
            String out = new String(p.getInputStream().readAllBytes()).trim();
            int code = p.waitFor();
            if (code != 0) {
                return null;
            }
            // reject empty / multi-line / non-hex-ish noise
            if (out.isEmpty() || out.contains("\n") || out.contains(" ")) {
                return null;
            }
            return out;
        } catch (IOException e) {
            return null;
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            return null;
        }
    }

    private static boolean gitDirty(Path root) {
        Path gitRoot = findGitRoot(root);
        if (gitRoot == null) {
            return false;
        }
        try {
            Process p = new ProcessBuilder("git", "-C", gitRoot.toString(), "status", "--porcelain")
                    .redirectErrorStream(true)
                    .start();
            boolean dirty = !new String(p.getInputStream().readAllBytes()).trim().isEmpty();
            int code = p.waitFor();
            return code == 0 && dirty;
        } catch (IOException e) {
            return false;
        } catch (InterruptedException e) {
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
