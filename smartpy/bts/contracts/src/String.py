import smartpy as sp

def split_btp_address(base):
    """
    Split the BTP Address format i.e. btp://1234.iconee/0x123456789
    into Network_address (1234.iconee) and Server_address (0x123456789)
    :param base: String base BTP Address format to be split
    :return: The resulting strings of Network_address and Server_address
    """
    sp.set_type(base, sp.TString)

    sep = sp.local("sep", "/")
    prev_idx = sp.local('prev_idx', 0)
    res = sp.local('res', [])
    sp.for idx in sp.range(0, sp.len(base)):
        sp.if sp.slice(base, idx, 1).open_some() == sep.value:
            res.value.push(sp.slice(base, prev_idx.value, sp.as_nat(idx - prev_idx.value)).open_some())
            prev_idx.value = idx + 1
    sp.if sp.len(base) > 0:
        res.value.push(sp.slice(base, prev_idx.value, sp.as_nat(sp.len(base) - prev_idx.value)).open_some())

    inverted_list = sp.local("my_list", res.value)
    last = sp.local("last", "")
    penultimate = sp.local("penultimate", "")

    with sp.match_cons(inverted_list.value) as l:
        last.value = l.head
        inverted_list.value = l.tail
    # with sp.else_():
    #     sp.failwith("Empty list")


    with sp.match_cons(inverted_list.value) as l:
        penultimate.value = l.head
    # with sp.else_():
    #     sp.failwith("Only one element")

    return sp.pair(penultimate.value, last.value)






