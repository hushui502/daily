package lambda;

import java.util.ArrayList;
import java.util.List;

public class OrderPattern {
    public interface Editor {
        public void save();
        public void open();
    }

    public interface Action {
        public void perform();
    }

    public class Save implements Action {
        private final Editor editor;
        public Save(Editor editor) {
            this.editor = editor;
        }
        @Override
        public void perform() {
            editor.save();
        }
    }

    public class Open implements Action {
        private final Editor editor;

        public Open(Editor editor) {
            this.editor = editor;
        }

        @Override
        public void perform() {
            editor.open();
        }
    }

    public static class Macro {
        private final List<Action> actions;

        public Macro() {
            actions = new ArrayList<>();
        }

        public void record(Action action) {
            actions.add(action);
        }

        public void run() {
            actions.forEach(Action::perform);
        }
    }

}
