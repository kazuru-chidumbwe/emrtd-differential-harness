package org.emrtd.harness.jmrtd;

import net.sf.scuba.smartcards.CardService;
import org.jmrtd.PassportService;
import org.jmrtd.lds.CardAccessFile;
import org.jmrtd.lds.PACEInfo;
import org.jmrtd.lds.SecurityInfo;

/**
 * JMRTD ≥0.7 removed the single-arg {@code PassportService(CardService)} constructor.
 * Pin the same tranceive/blocksize defaults for all harness runners (Option A: Maven Central 0.8.6).
 * {@link CardAccessFile#getPACEInfos()} was also removed — filter {@link SecurityInfo}.
 */
final class PassportServices {
    private PassportServices() {}

    static PassportService open(CardService card) {
        return new PassportService(
                card,
                PassportService.NORMAL_MAX_TRANCEIVE_LENGTH,
                PassportService.DEFAULT_MAX_BLOCKSIZE,
                false,
                true);
    }

    static PACEInfo firstPaceInfo(CardAccessFile cardAccess) {
        for (SecurityInfo info : cardAccess.getSecurityInfos()) {
            if (info instanceof PACEInfo) {
                return (PACEInfo) info;
            }
        }
        throw new IllegalStateException("CardAccessFile has no PACEInfo");
    }
}
