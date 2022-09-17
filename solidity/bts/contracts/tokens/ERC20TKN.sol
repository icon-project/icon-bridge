// SPDX-License-Identifier: Apache-2.0

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
pragma solidity >=0.8.2;

import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

contract ERC20TKN is Initializable, ContextUpgradeable, ERC20Upgradeable {
    uint8 decimal;

    function initialize(string memory name, string memory symbol, uint8 _decimals, uint _initialSupply) initializer public {
        require(_decimals <= 77, "OverLimit");
        decimal = _decimals;
        __Context_init_unchained();
       __ERC20_init(name, symbol);
        _mint(msg.sender, _initialSupply* (10**uint8(_decimals)));
    }
    
    function decimals() override public view returns ( uint8 ) {
        return decimal;
    }
}