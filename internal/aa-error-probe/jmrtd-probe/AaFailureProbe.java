import java.lang.reflect.Proxy;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.interfaces.RSAPublicKey;
import net.sf.scuba.smartcards.APDUWrapper;
import net.sf.scuba.smartcards.CardServiceException;
import org.jmrtd.APDULevelAACapable;
import org.jmrtd.protocol.AAProtocol;
import org.jmrtd.protocol.AAResult;

/**
 * Forced-failure probe for JMRTD AAProtocol.doAA (0.8.6).
 *
 * Case 1: INTERNAL AUTHENTICATE throws CardServiceException (chip SW fail).
 * Case 2: INTERNAL AUTHENTICATE returns garbage bytes with no exception —
 *         doAA still returns AAResult (no signature verification in AAProtocol).
 */
public class AaFailureProbe {
  public static void main(String[] args) throws Exception {
    KeyPairGenerator kpg = KeyPairGenerator.getInstance("RSA");
    kpg.initialize(2048);
    KeyPair kp = kpg.generateKeyPair();
    RSAPublicKey pub = (RSAPublicKey) kp.getPublic();
    byte[] challenge = new byte[] {1, 2, 3, 4, 5, 6, 7, 8};

    System.out.println("=== PROBE J1: doAA when sendInternalAuthenticate throws ===");
    APDULevelAACapable failService =
        (APDULevelAACapable)
            Proxy.newProxyInstance(
                APDULevelAACapable.class.getClassLoader(),
                new Class<?>[] {APDULevelAACapable.class},
                (proxy, method, methodArgs) -> {
                  if ("sendInternalAuthenticate".equals(method.getName())) {
                    throw new CardServiceException("forced SW 6982", 0x6982);
                  }
                  return null;
                });
    AAProtocol protoFail = new AAProtocol(failService, null);
    try {
      AAResult r = protoFail.doAA(pub, "SHA-256", "SHA256withRSA", challenge);
      System.out.println("AAResult returned (unexpected): " + r);
      System.out.println("VERDICT_J1: ERROR NOT THROWN");
    } catch (CardServiceException e) {
      System.out.println("threw: " + e.getClass().getName() + ": " + e.getMessage());
      System.out.println("cause: " + (e.getCause() == null ? "null" : e.getCause().toString()));
      System.out.println("VERDICT_J1: CardServiceException / protocol exception propagates");
    } catch (Exception e) {
      System.out.println("threw: " + e.getClass().getName() + ": " + e.getMessage());
      System.out.println("VERDICT_J1: Exception propagates (" + e.getClass().getSimpleName() + ")");
    }

    System.out.println();
    System.out.println("=== PROBE J2: doAA with OK APDU but garbage signature bytes ===");
    APDULevelAACapable garbageService =
        (APDULevelAACapable)
            Proxy.newProxyInstance(
                APDULevelAACapable.class.getClassLoader(),
                new Class<?>[] {APDULevelAACapable.class},
                (proxy, method, methodArgs) -> {
                  if ("sendInternalAuthenticate".equals(method.getName())) {
                    return new byte[] {0x00, 0x01, 0x02, 0x03};
                  }
                  return null;
                });
    AAProtocol protoOk = new AAProtocol(garbageService, null);
    try {
      AAResult r = protoOk.doAA(pub, "SHA-256", "SHA256withRSA", challenge);
      System.out.println("AAResult class=" + r.getClass().getName());
      System.out.println("response length=" + (r.getResponse() == null ? -1 : r.getResponse().length));
      System.out.println("VERDICT_J2: doAA returns AAResult WITHOUT verifying signature");
    } catch (Exception e) {
      System.out.println("threw: " + e.getClass().getName() + ": " + e.getMessage());
      System.out.println("VERDICT_J2: unexpected throw");
    }

    System.out.println();
    System.out.println("=== PROBE J3: bad challenge length ===");
    try {
      protoOk.doAA(pub, "SHA-256", "SHA256withRSA", new byte[] {1, 2, 3});
      System.out.println("VERDICT_J3: no throw (unexpected)");
    } catch (IllegalArgumentException e) {
      System.out.println("threw IllegalArgumentException: " + e.getMessage());
      System.out.println("VERDICT_J3: bad challenge rejected as IllegalArgumentException");
    } catch (Exception e) {
      System.out.println("threw: " + e.getClass().getName() + ": " + e.getMessage());
      System.out.println(
          "cause: " + (e.getCause() == null ? "null" : e.getCause().getClass().getName() + ": " + e.getCause().getMessage()));
      System.out.println(
          "VERDICT_J3: bad challenge becomes CardServiceProtocolException (AAProtocol catch-all wraps IllegalArgumentException)");
    }
  }
}
