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
package foundation.icon.btp.bmv.types;

import foundation.icon.btp.bmv.lib.mpt.MPTException;
import foundation.icon.btp.bmv.lib.mpt.Trie;
import score.Context;
import score.ObjectReader;
import scorex.util.ArrayList;

import java.util.List;

public class ReceiptProof {

    final static String RLPn = "RLPn";

    private final int index;
    private final byte[] mptKey;
  /*  private final List<EventProof> eventProofs;
    private final List<ReceiptEventLog> eventLogs;
    private final List<byte[]> mptProofs;*/
    private final List<EventDataBTPMessage> events;

    public ReceiptProof(int index, byte[] mptKey,/* List<byte[]> mptProofs, List<EventProof> eventProofs, List<ReceiptEventLog> eventLogs, */List<EventDataBTPMessage> events) {
        this.index = index;
        this.mptKey = mptKey;
       /* this.mptProofs = mptProofs;
        this.eventProofs = eventProofs;
        this.eventLogs = eventLogs;*/
        this.events = events;
    }

    public static ReceiptProof fromBytes(byte[] serialized) {
        ObjectReader reader = Context.newByteArrayObjectReader(RLPn, serialized);
        reader.beginList();
        //Index
        int index = reader.readInt();

        //mptKey
        byte[] mptKey = new byte[]{-128};

        if(index > 0)
           mptKey = new byte[]{(byte)index};

        /*List<byte[]> mptProofs = new ArrayList<>();
        //mptProofs.add(reader.readNullable(byte[].class));
        ObjectReader mptProofReader = Context.newByteArrayObjectReader(RLPn, reader.readNullable(byte[].class));
        try {
            mptProofReader.beginList();
            while (reader.hasNext()) {
                mptProofs.add(mptProofReader.readByteArray());
            }
            mptProofReader.end();
        } catch (Exception e) {
            //TODO: check why last reader.hasNext = true, even where there is no data, hence the exception
        }

        //EventProofs
        List<EventProof> eventProofs = readEventProofs(reader);

        //Event Logs
        List<ReceiptEventLog> eventsFromProofs = new ArrayList<>();
        for (EventProof ef : eventProofs) {
            eventsFromProofs.add(ReceiptEventLog.fromBytes(ef.getProof()));
        }*/
        List<EventDataBTPMessage> eventsLogs = new ArrayList<>();

        ObjectReader eventLogReader = Context.newByteArrayObjectReader(RLPn, reader.readByteArray());
        eventLogReader.beginList();
        while(eventLogReader.hasNext()){
            eventsLogs.add(EventDataBTPMessage.fromRLPBytes(eventLogReader));
        }
        eventLogReader.end();
        return new ReceiptProof(index, mptKey,/* mptProofs, eventProofs, eventsFromProofs,*/ eventsLogs);
    }

    public static List<byte[]> readByteArrayListFromRLP(byte[] serialized) {
        if (serialized == null)
            return null;
        ObjectReader reader = Context.newByteArrayObjectReader(RLPn, serialized);
        reader.beginList();
        List<byte[]> lists = new ArrayList<>();
        if (!reader.hasNext())
            return lists;

        while (reader.hasNext()) {
            lists.add(reader.readByteArray());
        }
        reader.end();

        return lists;
    }

    public static List<EventProof> readEventProofs(ObjectReader reader) {
        List<EventProof> eventProofs = new ArrayList<>();

        reader.beginList();
        while (reader.hasNext()) {
            reader.beginList();
            int index = reader.readInt();
            byte[] proof = reader.readNullable(byte[].class);
            eventProofs.add(new EventProof(index, proof));
            reader.end();
        }
        reader.end();

        return eventProofs;
    }

    public int getIndex() {
        return index;
    }

    public byte[] getMptKey() {
        return mptKey;
    }

    /*public List<byte[]> getMptProofs() {
        return mptProofs;
    }

    public List<EventProof> getEventProofs() {
        return eventProofs;
    }

    public List<ReceiptEventLog> getEventLogs() {
        return eventLogs;
    }*/

    public List<EventDataBTPMessage> getEvents() {
        return events;
    }
}
