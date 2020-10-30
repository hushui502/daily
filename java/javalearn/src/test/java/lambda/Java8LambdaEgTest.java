package lambda;

import org.junit.Test;

import java.util.*;
import java.util.function.BinaryOperator;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import static org.junit.Assert.*;

public class Java8LambdaEgTest {

    @Test
    public void testLambda() {
        List<String> list = Stream.of("a", "b", "c")
                .collect(Collectors.toList());

        assertEquals(Arrays.asList("a", "b", "c"), list);

        List<String> collected = Stream.of("a", "b", "c")
                .map(c -> c.toUpperCase())
                .collect(Collectors.toList());
        assertEquals(Arrays.asList("A", "B", "C"), collected);

        List<String> together = Stream.of(Arrays.asList("a", "b"), Arrays.asList("c", "d"))
                .flatMap(str -> str.stream())
                .collect(Collectors.toList());
        assertEquals(Arrays.asList("a", "b", "c", "d"), together);

        List<String> tracks = Arrays.asList("aaa", "bb", "c");
        String minTrack = tracks.stream()
                .min(Comparator.comparing(track -> track.length()))
                .get();

        assertEquals(tracks.get(2), minTrack);

        // reduce
        int count = Stream.of(1, 2, 3)
                .reduce(0, (x, y) -> x + y);
        assertEquals(6, count);

        BinaryOperator<Integer> accumulator = (x, y) -> x + y;
        int count1 = accumulator.apply(
                accumulator.apply(
                        accumulator.apply(0, 1),
                        2
                ),
                3
        );
        assertEquals(6, count1);
    }

    @Test
    public void TestFunc() {
        List<String> musics = Arrays.asList("hello", "beijing", "love");
        IntSummaryStatistics musicStatus = musics.stream()
                .mapToInt(music -> music.length())
                .summaryStatistics();
        assertEquals(4, musicStatus.getMin());
        musicStatus.getMax();
        musicStatus.getAverage();
        musicStatus.getSum();
    }

    @Test
    public void TestparentDefault() {
        Java8LambdaEg.Parent parent = new Java8LambdaEg.ParentImpl();
        parent.welcome();
    }

    @Test
    public void TestOptional() {
        Optional<String> a = Optional.of("a");
        assertEquals("a", a.get());
        assertTrue(a.isPresent());

        Optional emptyOptional = Optional.empty();
        Optional alsoEmpty = Optional.ofNullable(null);
        assertFalse(emptyOptional.isPresent());
        assertEquals("a", emptyOptional.orElse("a"));
        assertEquals("a", emptyOptional.orElseGet(() -> "a"));
    }

    @Test
    public void TestPartition() {

    }
}