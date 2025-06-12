const hre = require("hardhat");
const { upgrades } = require("hardhat");

async function main() {
    const BeggingContract = await hre.ethers.getContractFactory("BeggingContract");

    // 部署代理合约，传入初始化参数（例如 7 天）
    const durationInDays = 7;
    const beggingContract = await upgrades.deployProxy(BeggingContract, [durationInDays], {
        initializer: 'initialize',
        kind: 'transparent' // 指定使用透明代理
    });

    await beggingContract.waitForDeployment();

    console.log("BeggingContract (Transparent) deployed to:", await beggingContract.getAddress());
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});