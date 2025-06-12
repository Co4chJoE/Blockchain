const hre = require("hardhat");
const { upgrades } = require("hardhat");

async function main() {
    const BeggingContractV2 = await hre.ethers.getContractFactory("BeggingContractV2");
    console.log("Upgrading BeggingContract to V2...");


    const proxyAddress = "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9"; // 代理合约地址
    const upgradedProxy  = await upgrades.upgradeProxy(proxyAddress, BeggingContractV2);
    // await upgradedProxy.wait(); // 等待升级完成
    
    const upgradedImplementationAddress = await upgrades.erc1967.getImplementationAddress(proxyAddress);

    console.log("✅ 升级成功！");
    console.log("🔗 代理合约地址:", proxyAddress);
    console.log("🆕 实现合约地址:", upgradedImplementationAddress);
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});