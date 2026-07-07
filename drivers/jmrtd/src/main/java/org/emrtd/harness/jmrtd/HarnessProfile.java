package org.emrtd.harness.jmrtd;

import com.google.gson.Gson;
import com.google.gson.annotations.SerializedName;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;

public final class HarnessProfile {
    public String id;
    public String mechanism;
    public String condition;
    public Mrz mrz;
    @SerializedName("card_access_hex")
    public String cardAccessHex;
    public Injection injection;

    public static final class Mrz {
        @SerializedName("document_number")
        public String documentNumber;
        @SerializedName("date_of_birth")
        public String dateOfBirth;
        @SerializedName("date_of_expiry")
        public String dateOfExpiry;
    }

    public static final class Injection {
        @SerializedName("pace_sw")
        public String paceSw;
    }

    public static HarnessProfile load(Path path) throws IOException {
        String raw = Files.readString(path);
        HarnessProfile profile = new Gson().fromJson(raw, HarnessProfile.class);
        if (profile == null || profile.id == null || profile.mrz == null || profile.cardAccessHex == null) {
            throw new IOException("profile missing required fields");
        }
        return profile;
    }
}
