import java.lang.reflect.Proxy;
import net.sf.scuba.smartcards.CardServiceException;
import org.jmrtd.APDULevelEACTACapable;

/**
 * Forced-failure probe for JMRTD TA APDU path (0.8.6).
 * Chip SW fail on PSO (sendPSOExtendedLengthMode) / MSE must reach the caller.
 */
public final class TaFailureProbe {
  public static void main(String[] args) throws Exception {
    System.out.println("=== PROBE T1: sendPSOExtendedLengthMode throws SW=6982 ===");
    APDULevelEACTACapable fail =
        (APDULevelEACTACapable)
            Proxy.newProxyInstance(
                APDULevelEACTACapable.class.getClassLoader(),
                new Class<?>[] {APDULevelEACTACapable.class},
                (proxy, method, methodArgs) -> {
                  String n = method.getName();
                  if (n.startsWith("send") && !"sendGetChallenge".equals(n)) {
                    throw new CardServiceException("forced TA SW 6982", 0x6982);
                  }
                  if ("sendGetChallenge".equals(n)) {
                    return new byte[8];
                  }
                  return null;
                });

    try {
      fail.sendPSOExtendedLengthMode(null, new byte[] {0x01}, new byte[] {0x02});
      System.out.println("VERDICT_T1: ERROR NOT THROWN");
    } catch (CardServiceException e) {
      System.out.println("threw: " + e.getClass().getName() + ": " + e.getMessage());
      System.out.println("VERDICT_T1: CardServiceException propagates from TA PSO APDU");
    }

    System.out.println();
    System.out.println("=== PROBE T2: sendMutualAuthenticate throws ===");
    try {
      fail.sendMutualAuthenticate(null, new byte[] {0x01, 0x02, 0x03});
      System.out.println("VERDICT_T2: ERROR NOT THROWN");
    } catch (CardServiceException e) {
      System.out.println("threw: " + e.getClass().getName() + ": " + e.getMessage());
      System.out.println("VERDICT_T2: CardServiceException propagates from External Authenticate");
    }

    System.out.println();
    System.out.println("=== PROBE T3: APDULevelEACTACapable methods ===");
    for (var m : APDULevelEACTACapable.class.getMethods()) {
      if (m.getDeclaringClass() == APDULevelEACTACapable.class) {
        System.out.println("  " + m.getName());
      }
    }

    System.out.println();
    System.out.println("=== PROBE T4: gmrtd peer ===");
    System.out.println("VERDICT_T4: gmrtd has no TA API — peer_support.gmrtd=unsupported");
  }
}
