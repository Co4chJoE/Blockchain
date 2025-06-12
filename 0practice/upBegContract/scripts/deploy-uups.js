const hre = require("hardhat");
const { upgrades } = require("hardhat");

async function main() {
    const BeggingContract = await hre.ethers.getContractFactory("BeggingContract");

    // 使用 deployProxy 创建代理合约
    const proxy = await upgrades.deployProxy(BeggingContract, [7], {
        initializer: "initialize",
    }); 
    await proxy.waitForDeployment(); // 等待部署完成
    const proxyAddress = await proxy.getAddress();
    console.log("✅ UUPS 代理合约部署到:", proxyAddress);
}

main().catch((error) => {
  console.error("❌ 部署失败:", error);
  process.exitCode = 1;
});