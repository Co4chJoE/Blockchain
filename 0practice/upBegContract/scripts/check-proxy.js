const { ethers } = require("hardhat");
const { getImplementationAddress } = require("@openzeppelin/upgrades-core");

async function main() {
  const proxyAddress = "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9"; // 替换为你刚部署的地址

  try {
    const implAddress = await getImplementationAddress(ethers.provider, proxyAddress);
    console.log("✅ 这是一个代理合约，实现合约地址为:", implAddress);
  } catch (error) {
    console.error("❌ 这不是一个代理合约:", error.message);
  }
}

main().catch(console.error);