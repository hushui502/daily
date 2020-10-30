package lambda;

import java.awt.event.ActionListener;
import java.util.*;
import java.util.function.BinaryOperator;
import java.util.function.Predicate;
import java.util.stream.Collectors;
import java.util.stream.IntStream;
import java.util.stream.Stream;

public class Java8LambdaEg {
    Runnable runnable = () -> System.out.println("runnable");

    ActionListener actionListener = event -> System.out.println("action listener");

    Runnable multiStatement = () -> {
        System.out.println("hello");
        System.out.println("world");
    };

    BinaryOperator<Long> add = (x, y) -> x + y;
    BinaryOperator<Long> addExplicit = (Long x, Long y) -> x + y;

    // Function interface
    Predicate<Integer> atLeast5 = x -> x > 5;

    BinaryOperator<Float> add1 = (x, y) -> x + y;

    // stream


    public static void main(String[] args) {
        List<String> artList = new ArrayList<>();
        artList.add("hufan");
        artList.add("libai");
        Iterator<String> iterator = artList.iterator();

        while (iterator.hasNext()) {
            String res = iterator.next();
            System.out.println(res);
        }

        long count = artList.stream()
                .count();
        System.out.println(count);


        List<String> collected = Stream.of("a", "b", "c")
                .collect(Collectors.toList());
    }

    public long countPrimes(int upTo) {
        return IntStream.range(1, upTo)
                .filter(this::isPrime)
                .count();
    }

    private boolean isPrime(int number) {
        return IntStream.range(2, number)
                .allMatch(x -> (number % x) != 0);
    }


    public interface Parent {
        public void message(String body);

        public default void welcome() {
            message("hi");
        }

        public String getLastMessage();
    }

    public static class ParentImpl implements Parent {
        @Override
        public void message(String body) {

        }

        @Override
        public String getLastMessage() {
            return null;
        }
    }

    private int addIntegers(List<Integer> values) {
        return values.parallelStream()
                .mapToInt(i -> i)
                .sum();
    }

    private double[] parallelInitialize(int size) {
        double[] values = new double[size];
        Arrays.parallelSetAll(values, i -> i);
        return values;
    }

    public static double[] simpleMovingAverage(double[] values, int n) {
        double[] sums = Arrays.copyOf(values, values.length);
        Arrays.parallelPrefix(sums, Double::sum);
        int start = n - 1;
        return IntStream.range(start, sums.length)
                .mapToDouble(i -> {
                    double prefix = i == start ? 0 : sums[i-n];
                    return (sums[i] - prefix) / n;
                })
                .toArray();
    }

    public static List<String> allToUpperCase(List<String> words) {
        return words.stream()
                .map(string -> string.toUpperCase())
                .collect(Collectors.<String>toList());
    }

}
