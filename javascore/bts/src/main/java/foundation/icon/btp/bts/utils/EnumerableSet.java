package foundation.icon.btp.bts.utils;


import java.util.List;
import score.ArrayDB;
import score.Context;
import score.DictDB;
import scorex.util.ArrayList;

public class EnumerableSet<V> {

    private final ArrayDB<V> entries;
    private final DictDB<V, Integer> indexes;

    public EnumerableSet(String varKey, Class<V> valueClass) {
        // array of valueClass
        this.entries = Context.newArrayDB(varKey + "_es_entries", valueClass);
        // value => array index
        this.indexes = Context.newDictDB(varKey + "_es_indexes", Integer.class);
    }

    public int length() {
        return entries.size();
    }

    public V at(int index) {
        return entries.get(index);
    }

    public boolean contains(V value) {
        return indexes.get(value) != null;
    }

    public Integer indexOf(V value) {
        // returns null if value doesn't exist
        Integer result = indexes.get(value);
        if (result != null) {
            return result - 1;
        }
        return null;
    }

    public void add(V value) {
        if (!contains(value)) {
            // add new value
            entries.add(value);
            indexes.set(value, entries.size());
        }
    }

    public V remove(V value) {
        Integer valueIndex = indexOf(value);

        if (valueIndex != null) {
            int lastIndex = entries.size() - 1;
            V lastValue = entries.pop();
            indexes.set(value, null);
            if (lastIndex != valueIndex) {
                entries.set(valueIndex, lastValue);
                indexes.set(lastValue, valueIndex + 1);
                return lastValue;
            }
        }
        return null;
    }

    public List<V> range(int start, int end) {
        List<V> result = new ArrayList<>();
        int _end = Math.min(end, length());

        for (int i = start; i < _end; i++) {
            result.add(at(i));
        }
        return result;
    }
}
