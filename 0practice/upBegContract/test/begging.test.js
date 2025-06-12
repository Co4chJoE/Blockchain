const { expect } = require("chai");
const { ethers, upgrades } = require("hardhat");

describe("BeggingContract", function () {
  let beggingContract, owner, donor;

  beforeEach(async function () {
    const [ownerSigner, donorSigner] = await ethers.getSigners();
    owner = ownerSigner;
    donor = donorSigner;

    const BeggingContract = await ethers.getContractFactory("BeggingContract");
    beggingContract = await upgrades.deployProxy(BeggingContract, [7], { initializer: 'initialize' });
    await beggingContract.waitForDeployment();

    // 重新绑定完整合约 ABI，确保 getDonors() 可用(added)
    const contractAddress = await beggingContract.getAddress();
    beggingContract = await ethers.getContractAt("BeggingContract", contractAddress);
  });

  it("Should allow donations and track them correctly", async function () {
    const amount = ethers.parseEther("1");

    const contractAddress = await beggingContract.getAddress(); // 先获取地址(added)
    await donor.sendTransaction({ to: beggingContract.getAddress(), value: amount });

    expect(await beggingContract.getDonation(donor.address)).to.equal(amount);
  });

  it("Should allow owner to withdraw funds", async function () {
    const amount = ethers.parseEther("1");
    await donor.sendTransaction({ to: beggingContract.getAddress(), value: amount });

    const initialBalance = await ethers.provider.getBalance(owner.address);
    await beggingContract.connect(owner).withdraw();

    const finalBalance = await ethers.provider.getBalance(owner.address);
    expect(finalBalance - initialBalance).to.be.closeTo(amount, ethers.parseEther("0.01"));
  });

  //track donations(added in V2)
  it("Should track all donors correctly", async function () {

       //  使用 await 获取地址
        const contractAddress = await beggingContract.getAddress();
        // 检查初始捐赠者列表为空
        let donors = await beggingContract.getDonors();
        expect(donors.length).to.equal(0);

        // 第一个捐赠者捐款
        // await donor.sendTransaction({ to: contractAddress, value: ethers.parseEther("1") });
        await beggingContract.connect(donor).donate({ value: ethers.parseEther("1") });

        // 检查捐赠者列表包含第一个捐赠者
        donors = await beggingContract.getDonors();
        expect(donors).to.include(donor.address);
        expect(donors.length).to.equal(1);

        // 第二个捐赠者捐款（使用 owner 账户）
        // await owner.sendTransaction({ to: contractAddress, value: ethers.parseEther("1") });
        await beggingContract.connect(owner).donate({ value: ethers.parseEther("1") });

        // 检查捐赠者列表包含两个捐赠者
        donors = await beggingContract.getDonors();
        expect(donors).to.include(donor.address);
        expect(donors).to.include(owner.address);
        expect(donors.length).to.equal(2);
    });




});