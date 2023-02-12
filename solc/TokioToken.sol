// SPDX-License-Identifier: MIT

pragma solidity >=0.8.0 <=0.9.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract TokioToken is ERC20, Ownable {
    constructor() ERC20("TokioToken", "TKT") {
        _mint(msg.sender, 99*10**18);
    }

    function mint (address to, uint256 amount) public onlyOwner{
        _mint(to, amount);
    }
}
