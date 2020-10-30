package lambda;

public class LambdaTargetType {

    @FunctionalInterface
    interface A {
        void b();
    }

    @FunctionalInterface
    interface B {
        void b();
    }

    class UseAB {
        void use(A a) {
            System.out.println("Use A");
        }

        void use(B b) {
            System.out.println("Use B");
        }
    }

    void targetType() {
        UseAB useAB = new UseAB();
        A a = () -> System.out.println("Use");
        useAB.use(a);
    }
}
