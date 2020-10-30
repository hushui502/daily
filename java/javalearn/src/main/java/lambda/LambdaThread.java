package lambda;

import javax.swing.*;
import java.util.*;
import java.util.stream.Collectors;
import java.util.stream.Stream;

public class LambdaThread {
    public static void main(String[] args) {
        String name = "nameA";
        new Thread(() -> System.out.println("new thread start")).start();
        new Thread(() -> System.out.println("name is " + name)).start();


        // array stream
        Arrays.stream(new String[]{"hello", "world"})
                .forEach(System.out::println);

        int sum = Arrays.stream(new int[]{1, 2, 3})
                .reduce((a, b) -> a + b)
                .getAsInt();
        System.out.println(sum);

        // collection stream
        List<String> list = new ArrayList<>();
        list.add("a");
        list.add("b");
        list.stream()
                .forEach(System.out::println);

        // flatMap
        Stream.of(1, 2, 3)
                .map(v -> v + 1)
                .flatMap(v -> Stream.of(v*5, v*10))
                .forEach(System.out::println);

        // takeWhile dropWhile
        Stream.of(1, 2, 3)
                .takeWhile(v -> v < 3)
                .dropWhile(v -> v < 2)
                .forEach(System.out::println);

        System.out.println("reduce");
        // reduce
        Stream.of(2, 2, 3)
                .reduce((v1, v2) -> v1 + v2)
                .ifPresent(System.out::println);

        int result1 = Stream.of(1, 2, 3, 4, 5)
                .reduce(1, (v1, v2) -> v1 * v2);
        System.out.println("result1 " + result1);

        int result2 = Stream.of(1, 2, 3, 4, 5)
                .parallel()
                .reduce(0, (v1, v2) -> v1 + v2);
        System.out.println("result2 " + result2);

        // groupBy
        Map<Character, List<String>> names = Stream.of("Alex", "Alice", "Bob", "Bid", "Divid")
                .collect(Collectors.groupingBy(v -> v.charAt(0)));
        System.out.println(names);

        // joining
        String str = Stream.of("a", "b", "c")
                .collect(Collectors.joining(", "));
        System.out.println(str);

        // math
        double avgLength = Stream.of("hello", "world")
                .collect(Collectors.averagingInt(String::length));
        System.out.println(avgLength);

        IntSummaryStatistics statistics = Stream.of("a", "b", "c")
                .collect(Collectors.summarizingInt(String::length));

        System.out.println(statistics.getAverage());
        System.out.println(statistics.getCount());

        // button ActionListener
        JButton button = new JButton();
        button.addActionListener(event -> System.out.println("button clicked"));
    }
}


