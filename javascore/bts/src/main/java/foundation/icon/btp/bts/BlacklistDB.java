/*
 * Copyright 2021 ICON Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package foundation.icon.btp.bts;

import java.util.List;
import score.ArrayDB;
import score.BranchDB;
import score.Context;
import score.DictDB;
import scorex.util.ArrayList;

public class BlacklistDB {

    BranchDB<String, DictDB<String, Integer>> blacklistIndex = Context.newBranchDB("isBlacklistedDB", Integer.class);
    BranchDB<String, ArrayDB<String >> blacklistedUsers = Context.newBranchDB("blackListedUsers", String.class);

    public BlacklistDB() {}

    public int length(String net) {
        var a = blacklistedUsers.at(net);
        if (a != null) {
            return blacklistedUsers.at(net).size();
        }
        return 0;
    }

    public String at(String net, int index) {
        var a =  blacklistedUsers.at(net);
        if (a != null) {
            return a.get(index);
        }
        return null;
    }

    private String lowercase(String user) {
        return user.trim().toLowerCase();
    }

    public Integer indexOf(String net, String user) {
        var a = blacklistIndex.at(net);
        if (a != null) {
            Integer result = blacklistIndex.at(net).get(user);
            if (result != null) {
                return result - 1;
            }
            return null;
        }
        return null;
    }

    public boolean contains(String net, String user) {
        var a = blacklistIndex.at(net);
        if (a != null) {
            return a.get(lowercase(user)) != null;
        }
        return false;
    }

    public void addToBlacklist(String net, String user) {
        user = lowercase(user);
        if (!contains(net, user)) {
            blacklistedUsers.at(net).add(user);
            int size = length(net);
            blacklistIndex.at(net).set(user, size);
        }
    }

    public String removeFromBlacklist(String net, String user) {
        user = lowercase(user);
        Integer valueIdx = indexOf(net, user);
        var netUsers = blacklistedUsers.at(net);
        var netIndex = blacklistIndex.at(net);
        if (valueIdx != null) {
            int lastIdx = length(net) - 1;
            String lastVal = netUsers.pop();
            netIndex.set(user, null);
            if (lastIdx != valueIdx) {
                netUsers.set(valueIdx, lastVal);
                netIndex.set(lastVal, valueIdx +1);
                return lastVal;
            }
        }
        return null;
    }

    public List<String> range(String net, int start, int end) {
        List<String> result = new ArrayList<>();
        int _end = Math.min(end, length(net));
        for (int i = start; i < _end; i++) {
            result.add(at(net, i));
        }
        return result;
    }
}