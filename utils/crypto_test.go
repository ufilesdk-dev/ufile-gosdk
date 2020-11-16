package utils

import (
	"bytes"
	"encoding/hex"
	"testing"
)

var aesCTRTests = []struct {
	key, plaintext string
}{
	{
		"11754cd72aec309bf52f7687212e8957",
		"",
	},
	{
		"ca47248ac0b6f8372a97ac43508308ed",
		"",
	},
	{
		"fbe3467cc254f81be8e78d765a2e6333",
		"",
	},
	{
		"8a7f9d80d08ad0bd5a20fb689c88f9fc",
		"",
	},
	{
		"051758e95ed4abb2cdc69bb454110e82",
		"",
	},
	{
		"77be63708971c4e240d1cb79e8d77feb",
		"",
	},
	{
		"7680c5d3ca6154758e510f4d25b98820",
		"",
	},
	{
		"7fddb57453c241d03efbed3ac44e371c",
		"d5de42b461646c255c87bd2962d3b9a2",
	},
	{
		"ab72c77b97cb5fe9a382d9fe81ffdbed",
		"007c5e5b3e59df24a7c355584fc1518d",
	},
	{
		"fe47fcce5fc32665d2ae399e4eec72ba",
		"7c0e88c88899a779228465074797cd4c2e1498d259b54390b85e3eef1c02df60e743f1b840382c4bccaf3bafb4ca8429bea063",
	},
	{
		"ec0c2ba17aa95cd6afffe949da9cc3a8",
		"b85b3753535b825cbe5f632c0b843c741351f18aa484281aebec2f45bb9eea2d79d987b764b9611f6c0f8641843d5d58f3a242",
	},
	{
		"2c1f21cf0f6fb3661943155c3e3d8492",
		"42f758836986954db44bf37c6ef5e4ac0adaf38f27252a1b82d02ea949c8a1a2dbc0d68b5615ba7c1220ff6510e259f06655d8",
	},
	{
		"d9f7d2411091f947b4d6f1e2d1f0fb2e",
		"73ed042327f70fe9c572a61545eda8b2a0c6e1d6c291ef19248e973aee6c312012f490c2c6f6166f4a59431e182663fcaea05a",
	},
	{
		"fe9bb47deb3a61e423c2231841cfd1fb",
		"f1cc3818e421876bb6b8bbd6c9",
	},
	{
		"6703df3701a7f54911ca72e24dca046a",
		"793cd125b0b84a043e3ac67717",
	},
	// These cases test non-standard nonce sizes.
	{
		"1672c3537afa82004c6b8a46f6f0d026",
		"",
	},
	{
		"9a4fea86a621a91ab371e492457796c0",
		"ca6131faf0ff210e4e693d6c31c109fc5b6f54224eb120f37de31dc59ec669b6",
	},
	{
		"d0f1f4defa1e8c08b4b26d576392027c",
		"7ab49b57ddf5f62c427950111c5c4f0d",
	},
	{
		"4a0c00a3d284dea9d4bf8b8dde86685e",
		"6d4bf87640a6a48a50d28797b7",
	},
	{
		"0e18a844ac5bf38e4cd72d9b0942e506",
		"67c6697351ff4aec29cdbaabf2fbe3467cc254f81be8e78d765a2e63339fc99a66320db73158a35a255d051758e95ed4abb2cdc69bb454110e827441213ddc8770e93ea141e1fc673e017e97eadc6b968f385c2aecb03bfb32af3c54ec18db5c021afe43fbfaaa3afb29d1e6053c7c9475d8be6189f95cbba8990f95b1ebf1b3",
	},
	{
		"1f6c3a3bc0542aabba4ef8f6c7169e73",
		"67c6697351ff4aec29cdbaabf2fbe3467cc254f81be8e78d765a2e63339fc99a66320db73158a35a255d051758e95ed4abb2cdc69bb454110e827441213ddc8770e93ea141e1fc673e017e97eadc6b968f385c2aecb03bfb32af3c54ec18db5c021afe43fbfaaa3afb29d1e6053c7c9475d8be6189f95cbba8990f95b1ebf1b305eff700e9a13ae5ca0bcbd0484764bd1f231ea81c7b64c514735ac55e4b79633b706424119e09dcaad4acf21b10af3b33cde3504847155cbb6f2219ba9b7df50be11a1c7f23f829f8a41b13b5ca4ee8983238e0794d3d34bc5f4e77facb6c05ac86212baa1a55a2be70b5733b045cd33694b3afe2f0e49e4f321549fd824ea90870d4b28a2954489a0abcd50e18a844ac5bf38e4cd72d9b0942e506c433afcda3847f2dadd47647de321cec4ac430f62023856cfbb20704f4ec0bb920ba86c33e05f1ecd96733b79950a3e314d3d934f75ea0f210a8f6059401beb4bc4478fa4969e623d01ada696a7e4c7e5125b34884533a94fb319990325744ee9bbce9e525cf08f5e9e25e5360aad2b2d085fa54d835e8d466826498d9a8877565705a8a3f62802944de7ca5894e5759d351adac869580ec17e485f18c0c66f17cc07cbb22fce466da610b63af62bc83b4692f3affaf271693ac071fb86d11342d8def4f89d4b66335c1c7e4248367d8ed9612ec453902d8e50af89d7709d1a596c1f41f",
	},
	{
		"0795d80bc7f40f4d41c280271a2e4f7f",
		"1ad4e74d127f935beee57cff920665babe7ce56227377afe570ba786193ded3412d4812453157f42fafc418c02a746c1232c234a639d49baa8f041c12e2ef540027764568ce49886e0d913e28059a3a485c6eee96337a30b28e4cd5612c2961539fa6bc5de034cbedc5fa15db844013e0bef276e27ca7a4faf47a5c1093bd643354108144454d221b3737e6cb87faac36ed131959babe44af2890cfcc4e23ffa24470e689ce0894f5407bb0c8665cff536008ad2ac6f1c9ef8289abd0bd9b72f21c597bda5210cf928c805af2dd4a464d52e36819d521f967bba5386930ab5b4cf4c71746d7e6e964673457348e9d71d170d9eb560bd4bdb779e610ba816bf776231ebd0af5966f5cdab6815944032ab4dd060ad8dab880549e910f1ffcf6862005432afad",
	},
	{
		"e2e001a36c60d2bf40d69ff5b2b1161ea218db263be16a4e",
		"adb034f3f4a7ca45e2993812d113a9821d50df151af978bccc6d3bc113e15bc0918fb385377dca1916022ce816d56a332649484043c0fc0f2d37d040182b00a9bbb42ef231f80b48fb3730110d9a4433e38c73264c703579a705b9c031b969ec6d98de9f90e9e78b21179c2eb1e061946cd4bbb844f031ecf6eaac27a4151311adf1b03eda97c9fbae66295f468af4b35faf6ba39f9d8f95873bbc2b51cf3dfec0ed3c9b850696336cc093b24a8765a936d14dd56edc6bf518272169f75e67b74ba452d0aae90416a997c8f31e2e9d54ffea296dc69462debc8347b3e1af6a2d53bdfdfda601134f98db42b609df0a08c9347590c8d86e845bb6373d65a26ab85f67b50569c85401a396b8ad76c2b53ff62bcfbf033e435ef47b9b591d05117c6dc681d68e",
	},
	{
		"5394e890d37ba55ec9d5f327f15680f6a63ef5279c79331643ad0af6d2623525",
		"8e63067cd15359f796b43c68f093f55fdf3589fc5f2fdfad5f9d156668a617f7091d73da71cdd207810e6f71a165d0809a597df9885ca6e8f9bb4e616166586b83cc45f49917fc1a256b8bc7d05c476ab5c4633e20092619c4747b26dad3915e9fd65238ee4e5213badeda8a3a22f5efe6582d0762532026c89b4ca26fdd000eb45347a2a199b55b7790e6b1b2dba19833ce9f9522c0bcea5b088ccae68dd99ae0203c81b9f1dd3181c3e2339e83ccd1526b67742b235e872bea5111772aab574ae7d904d9b6355a79178e179b5ae8edc54f61f172bf789ea9c9af21f45b783e4251421b077776808f04972a5e801723cf781442378ce0e0568f014aea7a882dcbcb48d342be53d1c2ebfb206b12443a8a587cc1e55ca23beca385d61d0d03e9d84cbc1b0a",
	},
}

func TestAESCTR(t *testing.T) {
	for i, test := range aesCTRTests {
		key, _ := hex.DecodeString(test.key)
		plaintext, _ := hex.DecodeString(test.plaintext)

		crypto, err := NewCrypto(key)
		if err != nil {
			t.Errorf("#%d: NewCrypto err:%s", i, err.Error())
			continue
		}

		ct1 := crypto.XOR(plaintext)
		if err != nil {
			t.Errorf("#%d: Encrypt err:%s", i, err.Error())
			continue
		}

		/*ct2, err := crypto.Encrypt(plaintext)
		if err != nil {
			t.Errorf("#%d: Encrypt : %s", i, err.Error())
			continue
		}

		if bytes.Equal(ct1, ct2) {
			t.Errorf("#%d: ciphertext's match: got %x vs %x", i, ct1, ct2)
			continue
		}*/

		plaintext2 := crypto.XOR(ct1)
		if err != nil {
			t.Errorf("#%d: Decrypt : %s", i, err.Error())
			continue
		}
		if !bytes.Equal(plaintext, plaintext2) {
			t.Errorf("#%d: plaintext's don't match: got %x vs %x", i, plaintext2, plaintext)
			continue
		}
	}
}

func benchmarkAESCTREncrypt(b *testing.B, buf []byte) {
	b.SetBytes(int64(len(buf)))

	var key [16]byte

	crypto, err := NewCrypto(key[:])
	if err != nil {
		b.Errorf("NewCrypto: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		crypto.XOR(buf)

	}
}

func benchmarkAESCTRDecrypt(b *testing.B, buf []byte) {
	b.SetBytes(int64(len(buf)))

	var key [16]byte

	crypto, err := NewCrypto(key[:])
	if err != nil {
		b.Errorf("NewCrypto: %v", err)
	}

	out := crypto.XOR(buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = crypto.XOR(out)

	}
}

func BenchmarkAESCTREncrypt1K(b *testing.B) {
	benchmarkAESCTREncrypt(b, make([]byte, 1024))
}

func BenchmarkAESCTRDecrypt1K(b *testing.B) {
	benchmarkAESCTRDecrypt(b, make([]byte, 1024))
}

func BenchmarkAESCTREncrypt8K(b *testing.B) {
	benchmarkAESCTREncrypt(b, make([]byte, 8*1024))
}

func BenchmarkAESCTRDecrypt8K(b *testing.B) {
	benchmarkAESCTRDecrypt(b, make([]byte, 8*1024))
}
