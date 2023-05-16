import smartpy as sp



def split_btp_address(base, prev_string, result_string, list_string, last_string, penultimate_string):
    """
    Split the BTP Address format i.e. btp://1234.iconee/0x123456789
    into Network_address (1234.iconee) and Server_address (0x123456789)
    :param prev_string: local variable name
    :param result_string: local variable name
    :param list_string: local variable name
    :param last_string: local variable name
    :param penultimate_string: local variable name
    :param base: String base BTP Address format to be split
    :return: The resulting strings of Network_address and Server_address
    """
    sp.set_type(base, sp.TString)

    # sep = sp.local("sep", "/")
    prev_idx = sp.local(prev_string, 0)
    result = sp.local(result_string, [])
    sp.for idx in sp.range(0, sp.len(base)):
        sp.if sp.slice(base, idx, 1).open_some() == "/":
            result.value.push(sp.slice(base, prev_idx.value, sp.as_nat(idx - prev_idx.value)).open_some())
            prev_idx.value = idx + 1
    sp.if sp.len(base) > 0:
        result.value.push(sp.slice(base, prev_idx.value, sp.as_nat(sp.len(base) - prev_idx.value)).open_some())

    inverted_list = sp.local(list_string, result.value)
    last = sp.local(last_string, "")
    penultimate = sp.local(penultimate_string, "")

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






