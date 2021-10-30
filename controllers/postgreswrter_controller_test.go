package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	demov1 "github.com/yashvardhan-kukreja/postgres-writer-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("PostgresWriter controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		PostgresWriterName      = "sample-student"
		PostgresWriterNamespace = "default"
	)

	Context("In the beginning", func() {
		It("Should be able to create new PostgresWriter resource", func() {
			By("By creating a new PostgresWriter resource")
			ctx := context.Background()
			postgresWriter := &demov1.PostgresWriter{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "demo.yash.com/v1",
					Kind:       "PostgresWriter",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      PostgresWriterName,
					Namespace: PostgresWriterNamespace,
				},
				Spec: demov1.PostgresWriterSpec{
					Table:   "students",
					Name:    "Alex",
					Age:     1000,
					Country: "India",
				},
			}
			time.Sleep(10 * time.Second)
			Expect(k8sClient.Create(ctx, postgresWriter)).Should(Succeed())
		})
	})
})
