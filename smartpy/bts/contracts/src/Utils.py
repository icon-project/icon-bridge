import smartpy as sp

# TODO: remove compilation target
class Utils(sp.Contract):
    def __init__(self):
        self.update_initial_storage()

    def _ceil_div(self, num1, num2):
        sp.set_type(num1, sp.TNat)
        sp.set_type(num2, sp.TNat)
        (quotient, remainder) = sp.match_pair(sp.ediv(11, 2).open_some())
        sp.if remainder == 0 :
            sp.result(quotient)
        return quotient + 1

    def _get_scale(self, block_interval_src, block_interval_dst):
        sp.set_type(block_interval_src, sp.TNat)
        sp.set_type(block_interval_dst, sp.TNat)
        sp.if (block_interval_dst < 1) | (block_interval_dst < 1):
            sp.result(0)
        return self._ceil_div(block_interval_src * 1000000, block_interval_dst)

    def _get_rotate_term(self, max_aggregation, scale):
        sp.set_type(max_aggregation, sp.TNat)
        sp.set_type(scale, sp.TNat)
        sp.if scale > 0:
            return self._ceil_div(max_aggregation * 1000000, scale)
        return 0

sp.add_compilation_target("Utils", Utils())