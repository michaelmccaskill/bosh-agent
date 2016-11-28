package crypto_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"

	. "github.com/cloudfoundry/bosh-utils/crypto"
)

var _ = Describe("digest", func() {
	Describe("Digest", func() {
		Describe("#Verify", func() {
			It("verifies the algo and sum are matching", func() {
				expectedDigest := NewDigest(DigestAlgorithmSHA1, "07e1306432667f916639d47481edc4f2ca456454")
				actualDigest := NewDigest(DigestAlgorithmSHA1, "07e1306432667f916639d47481edc4f2ca456454")

				err := expectedDigest.Verify(actualDigest)
				Expect(err).ToNot(HaveOccurred())
			})

			Context("mismatching algorithm, matching digest", func() {
				It("errors", func() {
					expectedDigest := NewDigest(DigestAlgorithmSHA1, "07e1306432667f916639d47481edc4f2ca456454")
					actualDigest := NewDigest(DigestAlgorithmSHA256, "07e1306432667f916639d47481edc4f2ca456454")

					err := expectedDigest.Verify(actualDigest)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal(`Expected sha1 algorithm but received sha256`))
				})
			})

			Context("matching algorithm, mismatching digest", func() {
				It("errors", func() {
					expectedDigest := NewDigest(DigestAlgorithmSHA1, "07e1306432667f916639d47481edc4f2ca456454")
					actualDigest := NewDigest(DigestAlgorithmSHA1, "b1e66f505465c28d705cf587b041a6506cfe749f")

					err := expectedDigest.Verify(actualDigest)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal(`Expected sha1 digest "07e1306432667f916639d47481edc4f2ca456454" but received "b1e66f505465c28d705cf587b041a6506cfe749f"`))
				})
			})
		})

		Describe("#String", func() {
			Context("sha1", func() {
				It("excludes algorithm", func() {
					digest := NewDigest("sha1", "07e1306432667f916639d47481edc4f2ca456454")
					Expect(digest.String()).To(Equal("07e1306432667f916639d47481edc4f2ca456454"))
				})
			})

			Context("sha256", func() {
				It("includes algorithm", func() {
					digest := NewDigest("sha256", "b1e66f505465c28d705cf587b041a6506cfe749f7aa4159d8a3f45cc53f1fb23")
					Expect(digest.String()).To(Equal("sha256:b1e66f505465c28d705cf587b041a6506cfe749f7aa4159d8a3f45cc53f1fb23"))
				})
			})

			Context("sha512", func() {
				It("includes algorithm", func() {
					digest := NewDigest("sha512", "6f06a0c6c3827d827145b077cd8c8b7a15c75eb2bed809569296e6502ef0872c8e7ef91307a6994fcd2be235d3c41e09bfe1b6023df45697d88111df4349d64a")
					Expect(digest.String()).To(Equal("sha512:6f06a0c6c3827d827145b077cd8c8b7a15c75eb2bed809569296e6502ef0872c8e7ef91307a6994fcd2be235d3c41e09bfe1b6023df45697d88111df4349d64a"))
				})
			})
		})
	})

	Describe("ParseDigestString", func() {
		Describe("sha1", func() {
			It("creates a digest", func() {
				digest, err := ParseDigestString("sha1:07e1306432667f916639d47481edc4f2ca456454")
				Expect(err).ToNot(HaveOccurred())
				Expect(digest.Algorithm()).To(Equal(DigestAlgorithmSHA1))
				Expect(digest.Digest()).To(Equal("07e1306432667f916639d47481edc4f2ca456454"))
			})
		})

		Describe("sha256", func() {
			It("creates a digest", func() {
				digest, err := ParseDigestString("sha256:b1e66f505465c28d705cf587b041a6506cfe749f7aa4159d8a3f45cc53f1fb23")
				Expect(err).ToNot(HaveOccurred())
				Expect(digest.Algorithm()).To(Equal(DigestAlgorithmSHA256))
				Expect(digest.Digest()).To(Equal("b1e66f505465c28d705cf587b041a6506cfe749f7aa4159d8a3f45cc53f1fb23"))
			})
		})

		Describe("sha512", func() {
			It("creates a digest", func() {
				digest, err := ParseDigestString("sha512:6f06a0c6c3827d827145b077cd8c8b7a15c75eb2bed809569296e6502ef0872c8e7ef91307a6994fcd2be235d3c41e09bfe1b6023df45697d88111df4349d64a")
				Expect(err).ToNot(HaveOccurred())
				Expect(digest.Algorithm()).To(Equal(DigestAlgorithmSHA512))
				Expect(digest.Digest()).To(Equal("6f06a0c6c3827d827145b077cd8c8b7a15c75eb2bed809569296e6502ef0872c8e7ef91307a6994fcd2be235d3c41e09bfe1b6023df45697d88111df4349d64a"))
			})
		})

		Describe("default", func() {
			It("creates a sha1 digest", func() {
				digest, err := ParseDigestString("07e1306432667f916639d47481edc4f2ca456454")
				Expect(err).ToNot(HaveOccurred())
				Expect(digest.Algorithm()).To(Equal(DigestAlgorithmSHA1))
				Expect(digest.Digest()).To(Equal("07e1306432667f916639d47481edc4f2ca456454"))
			})
		})

		Describe("unrecognized", func() {
			It("errors", func() {
				_, err := ParseDigestString("unrecognized:something")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Unrecognized digest algorithm: unrecognized"))
			})
		})
	})

	Describe("CreateHashFromAlgorithm", func() {
		data := []byte("the checksum of c1oudc0w is deterministic")

		Describe("sha1", func() {
			It("hashes", func() {
				hash, err := CreateHashFromAlgorithm("sha1")
				Expect(err).ToNot(HaveOccurred())

				hash.Write(data)
				Expect(fmt.Sprintf("%x", hash.Sum(nil))).To(Equal("07e1306432667f916639d47481edc4f2ca456454"))
			})
		})

		Describe("sha256", func() {
			It("hashes", func() {
				hash, err := CreateHashFromAlgorithm("sha256")
				Expect(err).ToNot(HaveOccurred())

				hash.Write(data)
				Expect(fmt.Sprintf("%x", hash.Sum(nil))).To(Equal("b1e66f505465c28d705cf587b041a6506cfe749f7aa4159d8a3f45cc53f1fb23"))
			})
		})

		Describe("sha512", func() {
			It("hashes", func() {
				hash, err := CreateHashFromAlgorithm("sha512")
				Expect(err).ToNot(HaveOccurred())

				hash.Write(data)
				Expect(fmt.Sprintf("%x", hash.Sum(nil))).To(Equal("6f06a0c6c3827d827145b077cd8c8b7a15c75eb2bed809569296e6502ef0872c8e7ef91307a6994fcd2be235d3c41e09bfe1b6023df45697d88111df4349d64a"))
			})
		})

		Describe("unrecognized", func() {
			It("errors", func() {
				_, err := CreateHashFromAlgorithm("unrecognized")
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
