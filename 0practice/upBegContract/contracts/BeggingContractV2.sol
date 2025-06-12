// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "./BeggingContract.sol";

contract BeggingContractV2 is BeggingContract {
    // 新增功能：增加一个方法来查询总捐赠金额
    function getTotalDonations() public view returns (uint) {
        uint total = 0;
        address[] memory donors = getDonors();
        for (uint i = 0; i < donors.length; i++) {
            total += donations[donors[i]];
        }
        return total;
    }

     function getDonors() public view override returns (address[] memory) {
        // 返回 donors 数组
        address[] memory donors = new address[](2);
        return donors;
    }
    // 添加捐款函数，确保将新捐赠者加入列表
    function donate() public payable override {
        require(msg.value > 0, "Donation amount must be greater than zero");

        if (donations[msg.sender] == 0) {
            donors.push(msg.sender);
        }

        (bool success, ) = address(this).call{value: msg.value}("");
        require(success, "Transfer failed");
    }

}