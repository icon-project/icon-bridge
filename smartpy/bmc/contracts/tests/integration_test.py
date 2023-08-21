import subprocess
import smartpy as sp

t_balance_of_request = sp.TRecord(owner=sp.TAddress, token_id=sp.TNat).layout(("owner", "token_id"))
t_balance_of_response = sp.TRecord(request=t_balance_of_request, balance=sp.TNat).layout(("request", "balance"))

path_bts = {"bts_core": [
    {"types": ['types = sp.io.import_script_from_url("file:../bts/contracts/src/Types.py")',
               'types = sp.io.import_script_from_url("file:./contracts/src/Types.py")']},
    {"FA2_contract": [
        'FA2_contract = sp.io.import_script_from_url("file:../bts/contracts/src/FA2_contract.py")',
        'FA2_contract = sp.io.import_script_from_url("file:./contracts/src/FA2_contract.py")']}
],
    "bts_periphery": [
        {"types": ['types = sp.io.import_script_from_url("file:../bts/contracts/src/Types.py")',
                   'types = sp.io.import_script_from_url("file:./contracts/src/Types.py")']},
        {"strings": ['strings = sp.io.import_script_from_url("file:../bts/contracts/src/String.py")',
                     'strings = sp.io.import_script_from_url("file:./contracts/src/String.py")']},
        {"rlp": ['rlp = sp.io.import_script_from_url("file:../bts/contracts/src/RLP_struct.py")',
                 'rlp = sp.io.import_script_from_url("file:./contracts/src/RLP_struct.py")']},

    ],
    "RLP_struct": [
        {"types": ['types = sp.io.import_script_from_url("file:../bts/contracts/src/Types.py")',
                   'types = sp.io.import_script_from_url("file:./contracts/src/Types.py")']},
    ]
}

path_bmc = {
    "RLP_struct": [
        {'_to_byte':
            [
                '        _to_byte = sp.view("encode_nat", self.data.helper, sp.as_nat(params.sn), t=sp.TBytes).open_some()',
                '        _to_byte = sp.view("to_byte", self.data.helper_parse_negative, params.sn,'
                ' t=sp.TBytes).open_some()'
            ]
        },
        {'_to_int':
             ['                    _to_int = sp.to_int(Utils2.Int.of_bytes(sn_in_bytes))',
              '                    _to_int = sp.view("to_int", self.data.helper_parse_negative, sn_in_bytes, t=sp.TInt)'
              '.open_some()'
              ]
         }
    ],
    "bmc_periphery": [
        {'_self_address':
             ['        _self_address = sp.address("KT1T7kJERDbURCeZGBpTVq6yUR6pkWFJcfRG")',
              '        _self_address = sp.self_address'
              ]
         }
    ]
}

# in cas of new env: change above _self_address to bmc_periphery address of new env


def patch_file_path(file_name, old_value, new_value):
    subprocess.call("sed -i -e 's#.*" + old_value + " =.*#" + new_value + "#' " + file_name, shell=True)


def bts_core_contract_deploy_setup():
    for key, value in path_bts.items():
        for i in value:
            lis1 = []
            for x, y in i.items():
                lis1.append(x)
                # lis1.append(y)
                patch_file_path("../bts/contracts/src/" + key + ".py", lis1[0], y[0])


def bmc_periphery_contract_deploy_setup():
    for key, value in path_bmc.items():
        for i in value:
            lis1 = []
            for x, y in i.items():
                lis1.append(x)
                # lis1.append(y)
                patch_file_path("./contracts/src/" + key + ".py", lis1[0], y[0])


def tear_down_bts_changes():
    for key, value in path_bts.items():
        for i in value:
            lis1 = []
            for x, y in i.items():
                lis1.append(x)
                # lis1.append(y)
                patch_file_path("../bts/contracts/src/" + key + ".py", lis1[0], y[1])


def tear_down_bmc_changes():
    for key, value in path_bmc.items():
        for i in value:
            lis1 = []
            for x, y in i.items():
                lis1.append(x)
                # lis1.append(y)
                patch_file_path("./contracts/src/" + key + ".py", lis1[0], y[1])


# import changes in bts_core for testing
bts_core_contract_deploy_setup()
bmc_periphery_contract_deploy_setup()

BTSCore = sp.io.import_script_from_url("file:../bts/contracts/src/bts_core.py")
BTSOwnerManager = sp.io.import_script_from_url("file:../bts/contracts/src/bts_owner_manager.py")
BTSPeriphery = sp.io.import_script_from_url("file:../bts/contracts/src/bts_periphery.py")

FA2 = sp.io.import_script_from_url("https://legacy.smartpy.io/templates/fa2_lib.py")
BMCManagement = sp.io.import_script_from_url("file:./contracts/src/bmc_management.py")
BMCPeriphery = sp.io.import_script_from_url("file:./contracts/src/bmc_periphery.py")
BMCHelper = sp.io.import_script_from_url("file:./contracts/src/helper.py")
ParseAddress = sp.io.import_script_from_url("file:../bts/contracts/src/parse_address.py")
fa2_dummy_file = sp.io.import_script_from_url("file:./contracts/tests/fa2_dummy.py")

# revert the path changes made in bts_core for testing
tear_down_bts_changes()


@sp.add_test("IntegrationTest")
def test():
    sc = sp.test_scenario()

    # test account
    owner = sp.test_account("Owner")
    alice = sp.address("tz1ep7ffKsQCNdgnkPDCVnVkgbmFZP8mFN1G")
    bob = sp.address("tz1euHP1ntD4r3rv8BsE5pXpTRBnUFu69wYP")
    jack = sp.address("tz1UUvTndciyJJXuPHBEvWxuM9qgECMV31eA")
    sam = sp.address("tz1MrAHP91XLXJXBoB3WL52zQ8VDcnH5PeMp")
    relay = sp.test_account("Relay")
    helper_parse_neg_contract = sp.test_account("helper_parse_neg_contract")

    # in cas of new env: change icon_bmc_address and its block height of new environment
    icon_bmc_address = "cx674dbf2aae08b31ecb8174e755b2f0fa42a81298"
    icon_bmc_block_height = 10785336

    # deploy contracts

    helper_contract = deploy_helper_contract()
    sc += helper_contract

    bmc_management = deploy_bmc_management(owner.address, helper_contract.address)
    sc += bmc_management

    parse_address = deploy_parse_address()
    sc += parse_address

    bmc_periphery = deploy_bmc_periphery(bmc_management.address, helper_contract.address,
                                         helper_parse_neg_contract.address, parse_address.address)
    sc += bmc_periphery

    bts_owner_manager = deploy_bts_owner_manager_contract(owner.address)
    sc += bts_owner_manager

    bts_core = deploy_bts_core(bts_owner_manager.address)
    sc += bts_core

    bts_periphery = deploy_bts_periphery(bts_core.address, helper_contract.address,
                                         parse_address.address, bmc_periphery.address,
                                         owner.address)
    sc += bts_periphery

    fa2_dummy = fa2_dummy_file.SingleAssetToken(admin=owner.address,
                                                metadata=sp.utils.metadata_of_url(
                                                    "ipfs://example"),
                                                token_metadata=FA2.make_metadata(name="NativeWrappedCoin", decimals=6,
                                                                                 symbol="wTEZ"))
    sc += fa2_dummy

    fa2_dummy_second = fa2_dummy_file.SingleAssetToken(admin=owner.address,
                                                       metadata=sp.utils.metadata_of_url(
                                                           "ipfs://example"),
                                                       token_metadata=FA2.make_metadata(name="Dummy",
                                                                                        decimals=6,
                                                                                        symbol="PEPE"))
    sc += fa2_dummy_second

    # BMC_MANAGEMENT SETTERS

    # set pause status to false
    bmc_management.toggle_bridge_on().run(sender=owner.address)

    # set bmc periphery
    bmc_management.set_bmc_periphery(bmc_periphery.address).run(
        sender=owner.address)

    # set bmc_btp_address
    bmc_management.set_bmc_btp_address("NetXnHfVqm9iesp.tezos").run(
        sender=owner.address)

    # tear down changes after bmc btp address is set
    tear_down_bmc_changes()

    # add_service
    svc1 = sp.string("bts")
    bmc_management.add_service(sp.record(addr=bts_periphery.address, svc=svc1)).run(
        sender=owner.address)

    # add_route
    dst = "btp://0x7.icon/" + icon_bmc_address
    next_link = "btp://0x7.icon/" + icon_bmc_address
    bmc_management.add_route(sp.record(dst=dst, link=next_link)).run(
        sender=owner.address)

    # add_link
    bmc_management.add_link("btp://0x7.icon/" + icon_bmc_address).run(
        sender=owner.address)

    # set_link_rx_height
    bmc_management.set_link_rx_height(sp.record(height=icon_bmc_block_height,
                                                link="btp://0x7.icon/" + icon_bmc_address)).run(
        sender=owner.address)

    # add_relay
    bmc_management.add_relay(sp.record(link="btp://0x7.icon/" + icon_bmc_address,
                                       addr=sp.set([relay.address]))).run(
        sender=owner.address)

    # BTS_CORE SETTERS

    # set pause status to false
    bts_core.toggle_bridge_on().run(sender=owner.address)

    # update_bts_periphery
    bts_core.update_bts_periphery(bts_periphery.address).run(sender=owner.address)

    # set_fee_ratio
    bts_core.set_fee_ratio(name=sp.string("btp-NetXnHfVqm9iesp.tezos-XTZ"), fee_numerator=sp.nat(100),
                           fixed_fee=sp.nat(450)).run(sender=owner.address)

    prev = "btp://0x7.icon/" + icon_bmc_address

    # Tests
    # Scenario 1: Init message from relay
    msg_byte = sp.bytes(
        "0xf8e5f8e3b8e1f8df01b8d7f8d5f8d3b8406274703a2f2f4e6574586e486656716"
        "d39696573702e74657a6f732f4b543154376b4a45524462555243655a47427054567136"
        "79555236706b57464a6366524701b88ef88cb8396274703a2f2f3078372e69636f6e2f6378363734646"
        "26632616165303862333165636238313734653735356232663066613432613831323938b8406274703a2f2f4e6"
        "574586e486656716d39696573702e74657a6f732f4b543154376b4a45524462555243655a47427054567136795552"
        "36706b57464a6366524783626d630089c884496e697482c1c08400a4927b")
    bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg_byte)).run(sender=relay.address)

    # Scenario 2: Add to blacklist called from icon
    # add bob and tz1bPkYCh5rTTGL7DuPLB66J8zqnUD8cMRq1
    msg_byte = sp.bytes(
        "0xf9014df9014ab90147f9014401b9013bf90138f90135b8406274703a2f2f4e6574586e"
        "486656716d39696573702e74657a6f732f4b543154376b4a45524462555243655a47427054"
        "56713679555236706b57464a6366524702b8f0f8eeb8396274703a2f2f3078372e69636f6e2f63"
        "7836373464626632616165303862333165636238313734653735356232663066613432613831323938"
        "b8406274703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b543154376b4a45524462"
        "555243655a4742705456713679555236706b57464a636652478362747301b86af86803b865f86300f84"
        "aa4747a3165754850316e7444347233727638427345357058705452426e5546753639775950a4747a3162"
        "506b59436835725454474c3"
        "74475504c4236364a387a716e554438634d527131954e6574586e486656716d39696573702e74657a6f738400a4934b")
    bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg_byte)).run(sender=relay.address)
    # verify blacklisted address
    sc.verify_equal(bts_periphery.data.blacklist, {sp.address("tz1euHP1ntD4r3rv8BsE5pXpTRBnUFu69wYP"): True,
                                                   sp.address("tz1bPkYCh5rTTGL7DuPLB66J8zqnUD8cMRq1"): True})

    # Scenario 3: Remove from blacklist called from icon
    # remove bob
    msg_byte = sp.bytes(
        "0xf90127f90124b90121f9011e01b90115f90112f9010fb8406274703a2f2f4e6574586e486656716d39696573702e74657a6f732"
        "f4b543154376b4a45524462555243655a4742705456713679555236706b57464a6366524703b8caf8c8b8396274703"
        "a2f2f3078372e69636f6e2f637836373464626632616165303862333165636238313734653735356232663066613432613831323"
        "938b8406274703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b543154376b4a45524462555243655a47427054"
        "56713679555236706b57464a636652478362747302b844f84203b83ff83d01e5a4747a3165754850316e744434723372763842734"
        "5357058705452426e5546753639775950954e6574586e486656716d39696573702e74657a6f738400a49366")
    bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg_byte)).run(sender=relay.address)
    # verify one blacklisted address is removed
    sc.verify(bts_periphery.data.blacklist.get(sp.address("tz1bPkYCh5rTTGL7DuPLB66J8zqnUD8cMRq1")) == True)

    # Scenario 4: Transfer native coin from icon to tezos
    # transferred: ICX: 25*10**18
    # fee deducted on icon: 4550000000000000000
    # receiver address: bob

    # register icon native coin
    bts_core.register(
        name=sp.string("btp-0x7.icon-ICX"),
        fee_numerator=sp.nat(100),
        fixed_fee=sp.nat(4300000000000000000),
        addr=sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg"),
        token_metadata=sp.map({"token_metadata": sp.bytes("0x0dae11")}),
        metadata=sp.big_map({"metadata": sp.bytes("0x0dae11")})
    ).run(sender=owner.address)

    msg_byte = sp.bytes(
        "0xf90157f90154b90151f9014e01b90145f90142f9013fb8406274703a2f2f4e6574586e486656716d39"
        "696573702e74657a6f732f4b543154376b4a45524462555243655a4742705456713679555236706b57"
        "464a6366524704b8faf8f8b8396274703a2f2f3078372e69636f6e2f63783637346462663261616530"
        "3862333165636238313734653735356232663066613432613831323938b8406274703a2f2f4e6574586e"
        "486656716d39696573702e74657a6f732f4b543154376b4a45524462555243655a4742705456713679555"
        "236706b57464a636652478362747303b874f87200b86ff86daa6878396431383164313366343763356161"
        "65353535623730393831346336623232393738373937363139a4747a3165754850316e7444347233727638"
        "427345357058705"
        "452426e5546753639775950dcdb906274702d3078372e69636f6e2d49435889011bccfea6b8bd00008400a493c0")
    bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg_byte)).run(sender=relay.address)
    # verify native coin balance
    coin_address = bts_core.data.coins.get("btp-0x7.icon-ICX")
    user_balance = sp.view("get_balance_of", coin_address,
                           [sp.record(owner=bob, token_id=sp.nat(0))],
                           t=sp.TList(t_balance_of_response)).open_some("Invalid view")

    sc.verify_equal(user_balance, [sp.record(request=sp.record(owner=bob,
                                                               token_id=sp.nat(0)),
                                             balance=sp.nat(20450000000000000000))])

    # Scenario 5: Transfer ICON IRC2 coin bnUSD4 from icon to tezos without registering on tezos
    # transferred: bnUSD4: 50*10**18
    # receiver address: jack
    msg_byte = sp.bytes(
        "0xf9015af90157b90154f9015101b90148f90145f90142b8406274703a2f2f4e6574586e48665"
        "6716d39696573702e74657a6f732f4b543154376b4a45524462555243655a4742705456713679555"
        "236706b57464a6366524705b8fdf8fbb8396274703a2f2f3078372e69636f6e2f6378363734646266"
        "32616165303862333165636238313734653735356232663066613432613831323938b8406274703a2"
        "f2f4e6574586e486656716d39696573702e74657a6f732f4b543154376b4a45524462555243655a47"
        "42705456713679555236706b57464a636652478362747304b877f87500b872f870aa687839643138316"
        "431336634376335616165353535623730393831346336623232393738373937363139a4747a314d724"
        "148503931584c584a58426f4233574c35327a51385644636e48355"
        "0654d70dfde936274702d3078372e69636f6e2d626e555344348902b5e3af16b18800008400a494e9")
    bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg_byte)).run(sender=relay.address)
    # no changes happen on tezos so no need to assert

    # Scenario 6: Transfer ICON IRC2 coin bnUSD4 from icon to tezos
    # transferred: bnUSD4: 50*10**18
    # receiver address: jack
    # register icon coin bnUSD4
    bts_core.register(
        name=sp.string("btp-0x7.icon-bnUSD4"),
        fee_numerator=sp.nat(0),
        fixed_fee=sp.nat(0),
        addr=sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg"),
        token_metadata=sp.map({"token_metadata": sp.bytes("0x0dae11")}),
        metadata=sp.big_map({"metadata": sp.bytes("0x0dae11")})
    ).run(sender=owner.address)

    msg_byte = sp.bytes(
        "0xf9015af90157b90154f9015101b90148f90145f90142b8406274703a2f2f4e6574586e48665671"
        "6d39696573702e74657a6f732f4b543154376b4a45524462555243655a4742705456713679555236706b"
        "57464a63665247"
        "06b8fdf8fbb8396274703a2f2f3078372e69636f6e2f637836373464626632616165303862333165636238313"
        "734653735356232663066613432613831323938b8406274703a2f2f4e6574586e486656716d39696573702e74657a6"
        "f732f4b543154376b4a45524462555243655a4742705456713679555236706b57464a636652478362747305b877f"
        "87500b872f870aa68783964313831643133663437633561616535353562373039383134633662323239373837393736"
        "3139a4747a31555576546e646369794a4a587550484245765778754d39716745434d5633316541dfde936274702d3"
        "078372e69636f6e2d626e555344348902b5e3af16b18800008400a49561")
    bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg_byte)).run(sender=relay.address)
    # verify bnUSD4 coin balance
    coin_address = bts_core.data.coins.get("btp-0x7.icon-bnUSD4")
    user_balance = sp.view("get_balance_of", coin_address,
                           [sp.record(owner=jack, token_id=sp.nat(0))],
                           t=sp.TList(t_balance_of_response)).open_some("Invalid view")

    sc.verify_equal(user_balance, [sp.record(request=sp.record(owner=jack,
                                                               token_id=sp.nat(0)),
                                             balance=sp.nat(50000000000000000000))])

    # Scenario 7: Transfer batch from icon to tezos
    # transferred: icon native coin: 20*10**18 and bnUSD4:14*10**18
    # receiver address: alice
    msg_byte = sp.bytes(
        "0xf90179f90176b90173f9017001b90167f90164f90161b8406274703a2f2f4e6574586e486656716"
        "d39696573702e74657a6f732f4b543154376b4a45524462555243655a474270545671367955523670"
        "6b57464a6366524707b9011bf90118b8396274703a2f2f3078372e69636f6e2f63783637346462663"
        "2616165303862333165636238313734653735356232663066613432613831323938b8406274703a2f2"
        "f4e6574586e486656716d39696573702e74657a6f732f4b543154376b4a45524462555243655a47427"
        "05456713679555236706b57464a636652478362747306b894f89200b88ff88daa68783964313831643"
        "1336634376335616165353535623730393831346336623232393738373937363139a4747a316570376"
        "6664b7351434e64676e6b504443566e566b67626d465a50386d464e3147f83bdb906274702d3078372"
        "e69636f6e2d4943588900d71b0fe0"
        "a28e0000de936274702d3078372e69636f6e2d626e555344348900c249fdd3277800008400a4959d")
    bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg_byte)).run(sender=relay.address)
    # verify native coin balance
    coin_address = bts_core.data.coins.get("btp-0x7.icon-ICX")
    user_balance = sp.view("get_balance_of", coin_address,
                           [sp.record(owner=alice, token_id=sp.nat(0))],
                           t=sp.TList(t_balance_of_response)).open_some("Invalid view")

    sc.verify_equal(user_balance, [sp.record(request=sp.record(owner=alice,
                                                               token_id=sp.nat(0)),
                                             balance=sp.nat(15500000000000000000))])

    # verify bnUSD4 coin balance
    coin_address = bts_core.data.coins.get("btp-0x7.icon-bnUSD4")
    user_balance = sp.view("get_balance_of", coin_address,
                           [sp.record(owner=alice, token_id=sp.nat(0))],
                           t=sp.TList(t_balance_of_response)).open_some("Invalid view")

    sc.verify_equal(user_balance, [sp.record(request=sp.record(owner=alice,
                                                               token_id=sp.nat(0)),
                                             balance=sp.nat(14000000000000000000))])

    # Scenario 8: Transfer batch from icon to tezos
    # transferred 20 *10**18 bnUSD4
    # receiver address: alice
    msg_byte = sp.bytes(
        "0xf9015af90157b90154f9015101b90148f90145f90142b8406274703a2f2f4e6574586e486656716d3969657"
        "3702e74657a6f732f4b543154376b4a45524462555243655a4742705456713679555236706b57464a636652"
        "4708b8fdf8fbb8396274703a2f2f3078372e69636f6e2f63783637346462663261616530386233316563623"
        "8313734653735356232663066613432613831323938b8406274703a2f2f4e6574586e486656716d39696573"
        "702e74657a6f732f4b543154376b4a45524462555243655a4742705456713679555236706b57464a6366524"
        "78362747307b877f87500b872f870aa68783964313831643133663437633561616535353562373039383134"
        "6336623232393738373937363139a4747a3165703766664b7351434e64676e6b504443566e566b67626d465"
        "a50386d464e3147dfde936274702d3078372e69636f6e2d626e555344348901158e460913d000008400a495f3")
    bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg_byte)).run(sender=relay.address)

    # verify bnUSD4 coin balance
    # alice had 14*10**18 bnUSD4 initially
    coin_address = bts_core.data.coins.get("btp-0x7.icon-bnUSD4")
    user_balance = sp.view("get_balance_of", coin_address,
                           [sp.record(owner=alice, token_id=sp.nat(0))],
                           t=sp.TList(t_balance_of_response)).open_some("Invalid view")

    sc.verify_equal(user_balance, [sp.record(request=sp.record(owner=alice,
                                                               token_id=sp.nat(0)),
                                             balance=sp.nat(14000000000000000000 + 20000000000000000000))])

    # Scenario 9: Set token limit of bnUSD4 from icon
    # token limit of bnUSD4: 21 * 10**18
    msg_byte = sp.bytes(
        "0xf9011ef9011bb90118f9011501b9010cf90109f90106b8406274703a2f2f4e6574586e486656716d396965"
        "73702e74657a6f732f4b543154376b4a45524462555243655a4742705456713679555236706b57464a6366"
        "524709b8c1f8bfb8396274703a2f2f3078372e69636f6e2f63783637346462663261616530386233316563"
        "6238313734653735356232663066613432613831323938b8406274703a2f2f4e6574586e486656716d3969"
        "6573702e74657a6f732f4b543154376b4a45524462555243655a4742705456713679555236706b57464a63"
        "6652478362747308b83bf83904b7f6d4936274702d3078372e69636f6e2d626e55534434ca8901236efcbc"
        "bb340000954e6574586e486656716d39696573702e74657a6f738400a49624")
    bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg_byte)).run(sender=relay.address)
    # verify token limit
    sc.verify(bts_periphery.data.token_limit.get("btp-0x7.icon-bnUSD4") == sp.nat(21000000000000000000))

    # Scenario 10: Set token limit of btp-0x7.icon-bnUSD4 and btp-0x7.icon-ICX from icon
    # token limit of btp-0x7.icon-ICX: 43*10**18
    # token limit of btp-0x7.icon-bnUSD4: 32*10**18
    msg_byte = sp.bytes(
        "0xf9013bf90138b90135f9013201b90129f90126f90123b8406274703a2f2f4e6574586e486656716d39696573702e"
        "74657a6f732f4b543154376b4a45524462555243655a4742705456713679555236706b57464a636652470ab8def8"
        "dcb8396274703a2f2f3078372e69636f6e2f63783637346462663261616530386233316563623831373465373535"
        "6232663066613432613831323938b8406274703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b543"
        "154376b4a45524462555243655a4742705456713679555236706b57464a636652478362747309b858f85604b853f8"
        "51e5936274702d3078372e69636f6e2d626e55534434906274702d3078372e69636f6e2d494358d48901bc16d674e"
        "c800000890254beb02d1dcc0000954e6574586e486656716d39696573702e74657a6f738400a4965c")
    bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg_byte)).run(sender=relay.address)
    # verify token limits
    sc.verify(bts_periphery.data.token_limit.get("btp-0x7.icon-ICX") == sp.nat(43000000000000000000))
    sc.verify(bts_periphery.data.token_limit.get("btp-0x7.icon-bnUSD4") == sp.nat(32000000000000000000))

    # Scenario 11: Transfer btp-0x7.icon-bnUSD4 from icon to tezos
    # transferred btp-0x7.icon-bnUSD4: 32*10**18
    # receiver address: sam
    msg_byte = sp.bytes(
        "0xf9015af90157b90154f9015101b90148f90145f90142b8406274703a2f2f4e6574586e486656716d396965"
        "73702e74657a6f732f4b543154376b4a45524462555243655a4742705456713679555236706b57464a6366"
        "52470bb8fdf8fbb8396274703a2f2f3078372e69636f6e2f63783637346462663261616530386233316563"
        "6238313734653735356232663066613432613831323938b8406274703a2f2f4e6574586e486656716d3969"
        "6573702e74657a6f732f4b543154376b4a45524462555243655a4742705456713679555236706b57464a636"
        "65247836274730ab877f87500b872f870aa6878396431383164313366343763356161653535356237303938"
        "31346336623232393738373937363139a4747a314d724148503931584c584a58426f4233574c35327a51385"
        "644636e483550654d70dfde936274702d3078372e69636f6e2d626e555344348901bc16d674ec8000008400a496c8")
    bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg_byte)).run(sender=relay.address)
    # verify bnUSD4 coin balance
    coin_address = bts_core.data.coins.get("btp-0x7.icon-bnUSD4")
    user_balance = sp.view("get_balance_of", coin_address,
                           [sp.record(owner=sam, token_id=sp.nat(0))],
                           t=sp.TList(t_balance_of_response)).open_some("Invalid view")

    sc.verify_equal(user_balance, [sp.record(request=sp.record(owner=sam,
                                                               token_id=sp.nat(0)),
                                             balance=sp.nat(32000000000000000000))])

    # Tezos to icon scenarios

    # Scenario 12: Transfer native coin from tezos to icon
    # transferred btp-NetXnHfVqm9iesp.tezos-XTZ: 9000000
    # fee deducted on tezos: 90450

    bts_core.transfer_native_coin("btp://0x7.icon/hx9d181d13f47c5aae555b709814c6b22978797619").run(
        sender=alice, amount=sp.mutez(9000000))

    # relay msg for transfer end
    msg_byte = sp.bytes(
        "0xf8f2f8f0b8eef8ec01b8e4f8e2f8e0b8406274703a2f2f4e6574586e486656716d39696573702e74657"
        "a6f732f4b543154376b4a45524462555243655a4742705456713679555236706b57464a636652470cb89"
        "bf899b8396274703a2f2f3078372e69636f6e2f637836373464626632616165303862333165636238313"
        "734653735356232663066613432613831323938b8406274703a2f2f4e6574586e486656716d3969657370"
        "2e74657a6f732f4b543154376b4a45524462555243655a4742705456713679555236706b57464a63665247"
        "836274730196d50293d200905472616e7366657220537563636573738400a498d2")

    bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg_byte)).run(sender=relay.address)
    # verify aggregation fee
    # sc.verify(bts_core.data.aggregation_fee.get("btp-NetXnHfVqm9iesp.tezos-XTZ") == sp.nat(90450))

    # Scenario 13: Transfer wrapped token of btp-NetXnHfVqm9iesp.tezos-XTZ from icon to tezos
    # receiver address: tz1MrAHP91XLXJXBoB3WL52zQ8VDcnH5PeMp
    # received amount: 3*10**6

    msg_byte = sp.bytes("0xf9015ff9015cb90159f9015601b9014df9014af90147b8406274703a2f2f4e6"
                        "574586e486656716d39696573702e74657a6f732f4b543154376b4a45524462555"
                        "243655a4742705456713679555236706b57464a636652470db90101f8ffb839627"
                        "4703a2f2f3078372e69636f6e2f6378363734646266326161653038623331656362"
                        "38313734653735356232663066613432613831323938b8406274703a2f2f4e657458"
                        "6e486656716d39696573702e74657a6f732f4b543154376b4a45524462555243655a"
                        "4742705456713679555236706b57464a63665247836274730bb87bf87900b876f874"
                        "aa6878396431383164313366343763356161653535356237303938313463366232323"
                        "93738373937363139a4747a314d724148503931584c584a58426f4233574c35327a51"
                        "385644636e483550654d70e3e29d6274702d4e6574586e486656716d39696573702e7"
                        "4657a6f732d58545a832dc6c08400a49939")
    bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg_byte)).run(sender=relay.address)

    # Scenario 14: Transfer fa2 token from tezos to icon
    # mint fa2 dummy
    fa2_dummy.mint([sp.record(to_=alice, amount=sp.nat(1000000000))]).run(sender=owner.address)

    # add operator
    fa2_dummy.update_operators(
        [sp.variant("add_operator", sp.record(owner=alice, operator=bts_core.address, token_id=0))]).run(
        sender=alice)

    # register fa2
    bts_core.register(
        name=sp.string("btp-0x7.tezos-fa2"),
        fee_numerator=sp.nat(0),
        fixed_fee=sp.nat(0),
        addr=fa2_dummy.address,
        token_metadata=sp.map({"token_metadata": sp.bytes("0x0dae11")}),
        metadata=sp.big_map({"metadata": sp.bytes("0x0dae11")})
    ).run(sender=owner.address)

    # transfer fa2 token
    bts_core.transfer(sp.record(
        coin_name="btp-0x7.tezos-fa2", value=10000000,
        to="btp://0x7.icon/hx9d181d13f47c5aae555b709814c6b22978797619")).run(sender=alice)

    # relay msg for end of transfer fa2
    msg_byte = sp.bytes(
        "0xf8f2f8f0b8eef8ec01b8e4f8e2f8e0b8406274703a2f2f4e6574586e486656716d3969657370"
        "2e74657a6f732f4b543154376b4a45524462555243655a4742705456713679555236706b57464a6"
        "36652470eb89bf899b8396274703a2f2f3078372e69636f6e2f63783637346462663261616530386"
        "2333165636238313734653735356232663066613432613831323938b8406274703a2f2f4e6574586"
        "e486656716d39696573702e74657a6f732f4b543154376b4a45524462555243655a4742705456713"
        "679555236706b57464a63665247836274730296d50293d200905472616e7366657220537563636573738400a49c03")
    bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg_byte)).run(sender=relay.address)
    # verify aggregation fee
    sc.verify(bts_core.data.aggregation_fee.get("btp-0x7.tezos-fa2") == sp.nat(0))

    # Scenario 15: Transfer batch fa2 token and native token from tezos to icon
    bts_core.transfer_batch(sp.record(coin_names_values={"btp-0x7.tezos-fa2": 20000000},
                                      to="btp://0x7.icon/hx9d181d13f47c5aae555b709814c6b22978797619")).run(
        sender=alice, amount=sp.mutez(60000000))
    # relay msg for end of transfer batch
    msg_byte = sp.bytes(
        "0xf8f2f8f0b8eef8ec01b8e4f8e2f8e0b8406274703a2f2f4e6574586e486656716d3969657"
        "3702e74657a6f732f4b543154376b4a45524462555243655a4742705456713679555236706b5"
        "7464a636652470fb89bf899b8396274703a2f2f3078372e69636f6e2f637836373464626632"
        "616165303862333165636238313734653735356232663066613432613831323938b8406274703"
        "a2f2f4e6574586e486656716d39696573702e74657a6f732f4b543154376b4a455244625552436"
        "55a4742705456713679555236706b57464a63665247836274730396d50293d200905472616e736"
        "6657220537563636573738400a49e6e")
    bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg_byte)).run(sender=relay.address)
    # verify aggregation fee
    # existing fee of btp-NetXnHfVqm9iesp.tezos-XTZ: 90450
    sc.verify(bts_core.data.aggregation_fee.get("btp-0x7.tezos-fa2") == sp.nat(0))
    sc.verify(bts_core.data.aggregation_fee.get("btp-NetXnHfVqm9iesp.tezos-XTZ") == sp.nat(600450 + 90450))

    # Scenario 16: Transfer native coin from icon to tezos
    # receiving address: tz1g3pJZPifxhN49ukCZjdEQtyWgX2ERdfqP
    msg_byte = sp.bytes(
        "0xf90157f90154b90151f9014e01b90145f90142f9013fb8406274703a2f2f4e6574586e486"
        "656716d39696573702e74657a6f732f4b543154376b4a45524462555243655a4742705456713"
        "679555236706b57464a6366524710b8faf8f8b8396274703a2f2f3078372e69636f6e2f63783"
        "6373464626632616165303862333165636238313734653735356232663066613432613831323"
        "938b8406274703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b543154376b4a"
        "45524462555243655a4742705456713679555236706b57464a63665247836274730cb874f872"
        "00b86ff86daa68783964313831643133663437633561616535353562373039383134633662323"
        "2393738373937363139a4747a316733704a5a50696678684e3439756b435a6a644551747957675"
        "832455264667150dcdb906274702d3078372e69636f6e2d49435889011bccfea6b8bd00008400a49e9d")
    bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg_byte)).run(sender=relay.address)

    # Scenario 17: Transfer batch native coin and one fa2 tokens from tezos to icon
    bts_core.transfer_batch(sp.record(coin_names_values={"btp-0x7.tezos-fa2": 10000000},
                                      to="btp://0x7.icon/hx9d181d13f47c5aae555b709814c6b22978797619")).run(
        sender=alice, amount=sp.mutez(30000000))
    # relay msg for end of transfer batch
    msg_byte = sp.bytes(
        "0xf8f2f8f0b8eef8ec01b8e4f8e2f8e0b8406274703a2f2f4e6574586e486656716d3969"
        "6573702e74657a6f732f4b543154376b4a45524462555243655a4742705456713679555236"
        "706b57464a6366524711b89bf899b8396274703a2f2f3078372e69636f6e2f63783637346462"
        "6632616165303862333165636238313734653735356232663066613432613831323938b840627"
        "4703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b543154376b4a45524462555"
        "243655a4742705456713679555236706b57464a63665247836274730496d50293d200905472616"
        "e7366657220537563636573738400a49f1a")
    bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg_byte)).run(sender=relay.address)
    # verify aggregation fee
    # existing fee of btp-NetXnHfVqm9iesp.tezos-XTZ: 690900
    sc.verify(bts_core.data.aggregation_fee.get("btp-0x7.tezos-fa2") == sp.nat(0))
    sc.verify(bts_core.data.aggregation_fee.get("btp-NetXnHfVqm9iesp.tezos-XTZ") == sp.nat(690900 + 300450))

    # Scenario 18: Transfer fa2 token not registered on icon from tezos
    # mint fa2 dummy second
    fa2_dummy_second.mint([sp.record(to_=bob, amount=sp.nat(1000000000))]).run(sender=owner.address)

    # add operator
    fa2_dummy_second.update_operators(
        [sp.variant("add_operator", sp.record(owner=bob, operator=bts_core.address, token_id=0))]).run(
        sender=bob)

    # register fa2
    bts_core.register(
        name=sp.string("btp-0x7.tezos-fa2-second"),
        fee_numerator=sp.nat(0),
        fixed_fee=sp.nat(0),
        addr=fa2_dummy_second.address,
        token_metadata=sp.map({"token_metadata": sp.bytes("0x0dae11")}),
        metadata=sp.big_map({"metadata": sp.bytes("0x0dae11")})
    ).run(sender=owner.address)

    # transfer fa2 token
    bts_core.transfer(sp.record(
        coin_name="btp-0x7.tezos-fa2-second", value=10000000,
        to="btp://0x7.icon/hx9d181d13f47c5aae555b709814c6b22978797619")).run(sender=bob)

    user_balance_before = sp.view("get_balance_of", fa2_dummy_second.address,
                                  [sp.record(owner=bob, token_id=sp.nat(0))],
                                  t=sp.TList(t_balance_of_response)).open_some("Invalid view")
    # relay msg for end of transfer fa2
    msg_byte = sp.bytes(
        "0xf8e1f8dfb8ddf8db01b8d3f8d1f8cfb8406274703a2f2f4e6574586e486656"
        "716d39696573702e74657a6f732f4b543154376b4a45524462555243655a47427"
        "05456713679555236706b57464a6366524712b88af888b8396274703a2f2f3078"
        "372e69636f6e2f6378363734646266326161653038623331656362383137346537"
        "35356232663066613432613831323938b8406274703a2f2f4e6574586e48665671"
        "6d39696573702e74657a6f732f4b543154376b4a45524462555243655a47427054"
        "56713679555236706b57464a636652478362747381fb84c328f8008400a4b09d")
    bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg_byte)).run(sender=relay.address)
    user_balance_after = sp.view("get_balance_of", fa2_dummy_second.address,
                                  [sp.record(owner=bob, token_id=sp.nat(0))],
                                  t=sp.TList(t_balance_of_response)).open_some("Invalid view")

    sc.verify_equal(user_balance_before, user_balance_after)

    # sc.verify_equal(user_balance, [sp.record(request=sp.record(owner=sam,
    #                                                            token_id=sp.nat(0)),
    #                                          balance=sp.nat(32000000000000000000))])

    # these case cannot be tested in integration test due to limitation on tezos
    # # Scenario 12: Transferred btp-0x7.icon-ICX wrapped coin from tezos to icon
    # coin_address = bts_core.data.coins.get("btp-0x7.icon-ICX")
    # # contract = sp.contract(sp.TRecord(spender=sp.TAddress, amount=sp.TNat), coin_address, "set_allowance").open_some()
    #
    # # set allowance for bts_core
    # coin_address.set_allowance([sp.record(spender=bts_core.address,
    #                                       amount=sp.nat(1000000000000000000000))])
    # sc.verify(sp.view("get_allowance", coin_address, sp.record(spender=bts_core.address, owner=alice),
    #                   t=sp.TNat).open_some("Invalid view") == sp.nat(0))
    # # update operator
    # coin_address.update_operators(
    #     [sp.variant("add_operator", sp.record(owner=alice, operator=bts_core.address, token_id=0))])
    #
    # bts_balance_before = sp.view("get_balance_of", coin_address,
    #                              [sp.record(owner=bts_core.address, token_id=sp.nat(0))],
    #                              t=sp.TList(t_balance_of_response)).open_some("Invalid view")
    #
    # # transfer wrapped coin
    # bts_core.transfer(sp.record(coin_name="btp-0x7.icon-ICX", value=sp.nat(13000000000000000000),
    #                             to=" btp://0x7.icon/hx9d181d13f47c5aae555b709814c6b22978797619")).run(sender=alice)
    # # transfer end message from relay
    # msg_byte = sp.bytes(
    #     "0xf8f2f8f0b8eef8ec01b8e4f8e2f8e0b8406274703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b5431574b42484c"
    #     "4662674c38514a574779536635434b32684b3369434c454d395a69540cb89bf899b8396274703a2f2f3078372e69636f6e2f63783531"
    #     "6530626238353833396530653366666662346330313430616530663038336538393834363464b8406274703a2f2f4e6574586e486656"
    #     "716d39696573702e74657a6f732f4b5431574b42484c4662674c38514a574779536635434b32684b3369434c454d395a695483627473"
    #     "0196d50293d200905472616e7366657220537563636573738400a2b2de")
    # bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg_byte)).run(sender=relay.address)
    #
    # bts_balance_after = sp.view("get_balance_of", coin_address,
    #                             [sp.record(owner=bts_core.address, token_id=sp.nat(0))],
    #                             t=sp.TList(t_balance_of_response)).open_some("Invalid view")
    # verify burnt amount

    # Scenario 13: Transferred btp-0x7.icon-ICX wrapped coin from tezos to icon


def deploy_bmc_management(owner, helper):
    bmc_management = BMCManagement.BMCManagement(owner, helper)
    return bmc_management


def deploy_bmc_periphery(bmc_address, helper, helper_parse_neg_contract, parse):
    bmc_periphery = BMCPeriphery.BMCPreiphery(bmc_address, helper, helper_parse_neg_contract, parse)
    return bmc_periphery


def deploy_helper_contract():
    helper_contract = BMCHelper.Helper()
    return helper_contract


def deploy_parse_address():
    parse_address = ParseAddress.ParseAddress()
    return parse_address


def deploy_bts_core(bts_owner_manager_contract):
    bts_core = BTSCore.BTSCore(
        owner_manager=bts_owner_manager_contract,
        _native_coin_name="btp-NetXnHfVqm9iesp.tezos-XTZ",
        _fee_numerator=sp.nat(100),
        _fixed_fee=sp.nat(450)
    )
    return bts_core


def deploy_bts_owner_manager_contract(owner):
    bts_owner_manager_contract = BTSOwnerManager.BTSOwnerManager(owner)
    return bts_owner_manager_contract


def deploy_bts_periphery(core_address, helper, parse, bmc, owner):
    bts_periphery_contract = BTSPeriphery.BTSPeriphery(bmc_address=bmc, bts_core_address=core_address,
                                                       helper_contract=helper, parse_address=parse,
                                                       owner_address=owner,
                                                       native_coin_name="btp-NetXnHfVqm9iesp.tezos-XTZ")
    return bts_periphery_contract
