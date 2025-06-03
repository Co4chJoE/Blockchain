// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

contract BeggingContract {
    // 合约所有者地址
    address public owner;

    // 捐赠记录 mapping(address => uint)
    mapping(address => uint) public donations;

    // 捐赠事件
    event Donation(address indexed donor, uint amount, uint timestamp);

    // 提款事件
    event Withdrawal(address indexed admin, uint amount, uint timestamp);

    // 时间限制参数（可选）
    uint public startTime;
    uint public endTime;

    // 构造函数，设置合约所有者，并设定捐赠时间段
    constructor(uint durationInDays) {
        owner = msg.sender;
        startTime = block.timestamp;
        endTime = block.timestamp + (durationInDays * 1 days);
    }

    // 捐赠函数 payable
    function donate() public payable {
        // require(block.timestamp >= startTime && block.timestamp <= endTime, "The donation period has ended.");
        emit Donation(msg.sender, msg.value, block.timestamp);

        require(msg.value > 0, "The donation amount must be greater than 0.");

        donations[msg.sender] += msg.value;
        // emit Donation(msg.sender, msg.value);

    }
    // 提款函数 payable
    function withdraw() public onlyOwner {
        uint balance = address(this).balance;
        require(balance > 0, "There is no fund available for withdrawal in the contract.");

        payable(owner).transfer(balance);
        // emit Withdrawal(owner, balance);

        emit Withdrawal(owner, balance, block.timestamp);
    }

    // 查询某个地址的捐赠金额
    function getDonation(address donor) public view returns (uint) {
        return donations[donor];
    }

    receive() external payable {
        if (msg.value > 0) {
            donations[msg.sender] += msg.value;
            emit Donation(msg.sender, msg.value, block.timestamp);
        }
    }

    fallback() external payable {
        if (msg.value > 0) {
            donations[msg.sender] += msg.value;
            emit Donation(msg.sender, msg.value, block.timestamp);
        }
    }


    // 只有合约所有者才能调用的修饰符
    modifier onlyOwner {
        require(msg.sender == owner, "Only the owner of the contract has permission to perform this action.");
        _;
    }
}