// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

contract BeggingContract is OwnableUpgradeable, UUPSUpgradeable {
    // 合约所有者地址
    // address public owner;
    // uint256 public duration;
    // 捐赠记录 mapping(address => uint)
    mapping(address => uint) public donations;

    // 捐赠事件
    event Donation(address indexed donor, uint amount, uint timestamp);

    // 提款事件
    event Withdrawal(address indexed admin, uint amount, uint timestamp);


    // 记录捐赠者地址和去重映射(added)
    // donors 数组存储所有捐赠者地址
    address[] internal donors;
    mapping(address => bool) private _isDonor;

    // 时间限制参数（可选）
    uint public startTime;
    uint public endTime;

    // 构造函数，设置合约所有者，并设定捐赠时间段
    function initialize(uint durationInDays) initializer public {
        // owner = msg.sender;
        startTime = block.timestamp;
        endTime = block.timestamp + (durationInDays * 1 days);
        __Ownable_init(msg.sender); // 初始化权限控制
    }


    // 实现 UUPSUpgradeable 的授权升级逻辑(added)
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {
        // 只有 owner 可以触发升级
        // 可以在这里加入额外逻辑，比如暂停检查、版本控制等
    }

    // 捐赠函数 payable
    function donate() public payable virtual {
        // require(block.timestamp >= startTime && block.timestamp <= endTime, "The donation period has ended.");
        emit Donation(msg.sender, msg.value, block.timestamp);

        require(msg.value > 0, "The donation amount must be greater than 0.");

        donations[msg.sender] += msg.value;
        // emit Donation(msg.sender, msg.value);

        // 如果是第一次捐赠，加入 donors 列表(added)
        if (!_isDonor[msg.sender]) {
            donors.push(msg.sender);
            _isDonor[msg.sender] = true;
        }

        emit Donation(msg.sender, msg.value, block.timestamp);


    }

    // 为升级用
    function getDonors() public view virtual returns (address[] memory) {
        // 返回所有捐赠者地址
        return donors;
    }


    // 提款函数 payable
    function withdraw() public onlyOwner {
        uint balance = address(this).balance;
        require(balance > 0, "There is no fund available for withdrawal in the contract.");

        payable(owner()).transfer(balance);
        // emit Withdrawal(owner, balance);

        emit Withdrawal(owner(), balance, block.timestamp);
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
    // modifier onlyOwner {
    //     require(msg.sender == owner, "Only the owner of the contract has permission to perform this action.");
    //     _;
    // }
}


// 待实现的升级授权函数
// function _authorizeUpgrade(address) internal override onlyOwner {}