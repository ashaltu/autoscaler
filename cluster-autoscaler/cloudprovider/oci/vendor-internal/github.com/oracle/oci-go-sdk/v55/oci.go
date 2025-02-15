/*
This is the official Go SDK for Oracle Cloud Infrastructure

# Installation

Refer to https://github.com/oracle/oci-go-sdk/blob/master/README.md#installing for installation instructions.

# Configuration

Refer to https://github.com/oracle/oci-go-sdk/blob/master/README.md#configuring for configuration instructions.

# Quickstart

The following example shows how to get started with the SDK. The example belows creates an identityClient
struct with the default configuration. It then utilizes the identityClient to list availability domains and prints
them out to stdout

	import (
		"context"
		"fmt"

		"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/oci/vendor-internal/github.com/oracle/oci-go-sdk/v55/common"
		"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/oci/vendor-internal/github.com/oracle/oci-go-sdk/v55/identity"
	)

	func main() {
		c, err := identity.NewIdentityClientWithConfigurationProvider(common.DefaultConfigProvider())
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// The OCID of the tenancy containing the compartment.
		tenancyID, err := common.DefaultConfigProvider().TenancyOCID()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		request := identity.ListAvailabilityDomainsRequest{
			CompartmentId: &tenancyID,
		}

		r, err := c.ListAvailabilityDomains(context.Background(), request)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf("List of available domains: %v", r.Items)
		return
	}

More examples can be found in the SDK Github repo: https://github.com/oracle/oci-go-sdk/tree/master/example

# Optional fields in the SDK

Optional fields are represented with the `mandatory:"false"` tag on input structs. The SDK will omit all optional fields that are nil when making requests.
In the case of enum-type fields, the SDK will omit fields whose value is an empty string.

# Helper functions

The SDK uses pointers for primitive types in many input structs. To aid in the construction of such structs, the SDK provides
functions that return a pointer for a given value. For example:

	// Given the struct
	type CreateVcnDetails struct {

		// Example: `172.16.0.0/16`
		CidrBlock *string `mandatory:"true" json:"cidrBlock"`

		CompartmentId *string `mandatory:"true" json:"compartmentId"`

		DisplayName *string `mandatory:"false" json:"displayName"`

	}

	// We can use the helper functions to build the struct
	details := core.CreateVcnDetails{
		CidrBlock:     common.String("172.16.0.0/16"),
		CompartmentId: common.String("someOcid"),
		DisplayName:   common.String("myVcn"),
	}

# Customizing Requests

The SDK exposes functionality that allows the user to customize any http request before is sent to the service.

You can do so by setting the `Interceptor` field in any of the `Client` structs. For example:

	client, err := audit.NewAuditClientWithConfigurationProvider(common.DefaultConfigProvider())
	if err != nil {
		panic(err)
	}

	// This will add a header called "X-CustomHeader" to all request
	// performed with client
	client.Interceptor = func(request *http.Request) error {
		request.Header.Set("X-CustomHeader", "CustomValue")
		return nil
	}

The Interceptor closure gets called before the signing process, thus any changes done to the request will be properly
signed and submitted to the service.

# Signing Custom Requests

The SDK exposes a stand-alone signer that can be used to signing custom requests. Related code can be found here:
https://github.com/oracle/oci-go-sdk/blob/master/common/http_signer.go.

The example below shows how to create a default signer.

	client := http.Client{}
	var request http.Request
	request = ... // some custom request

	// Set the Date header
	request.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))

	// And a provider of cryptographic keys
	provider := common.DefaultConfigProvider()

	// Build the signer
	signer := common.DefaultSigner(provider)

	// Sign the request
	signer.Sign(&request)

	// Execute the request
	client.Do(request)

The signer also allows more granular control on the headers used for signing. For example:

	client := http.Client{}
	var request http.Request
	request = ... // some custom request

	// Set the Date header
	request.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))

	// Mandatory headers to be used in the sign process
	defaultGenericHeaders    = []string{"date", "(request-target)", "host"}

	// Optional headers
	optionalHeaders = []string{"content-length", "content-type", "x-content-sha256"}

	// A predicate that specifies when to use the optional signing headers
	optionalHeadersPredicate := func (r *http.Request) bool {
		return r.Method == http.MethodPost
	}

	// And a provider of cryptographic keys
	provider := common.DefaultConfigProvider()

	// Build the signer
	signer := common.RequestSigner(provider, defaultGenericHeaders, optionalHeaders, optionalHeadersPredicate)

	// Sign the request
	signer.Sign(&request)

	// Execute the request
	c.Do(request)

You can combine a custom signer with the exposed clients in the SDK.
This allows you to add custom signed headers to the request. Following is an example:

	//Create a provider of cryptographic keys
	provider := common.DefaultConfigProvider()

	//Create a client for the service you interested in
	c, _ := identity.NewIdentityClientWithConfigurationProvider(provider)

	//Define a custom header to be signed, and add it to the list of default headers
	customHeader := "opc-my-token"
	allHeaders := append(common.DefaultGenericHeaders(), customHeader)

	//Overwrite the signer in your client to sign the new slice of headers
	c.Signer = common.RequestSigner(provider, allHeaders, common.DefaultBodyHeaders())

	//Set the value of the header. This can be done with an Interceptor
	c.Interceptor = func(request *http.Request) error {
		request.Header.Add(customHeader, "customvalue")
		return nil
	}

	//Execute your operation as before
	c.ListGroups(..)

Bear in mind that some services have a white list of headers that it expects to be signed.
Therefore, adding an arbitrary header can result in authentications errors.
To see a runnable example, see https://github.com/oracle/oci-go-sdk/blob/master/example/example_identity_test.go

For more information on the signing algorithm refer to: https://docs.cloud.oracle.com/Content/API/Concepts/signingrequests.htm

# Polymorphic json requests and responses

Some operations accept or return polymorphic json objects. The SDK models such objects as interfaces. Further the SDK provides
structs that implement such interfaces. Thus, for all operations that expect interfaces as input, pass the struct in the SDK that satisfies
such interface. For example:

	c, err := identity.NewIdentityClientWithConfigurationProvider(common.DefaultConfigProvider())
	if err != nil {
		panic(err)
	}

	// The CreateIdentityProviderRequest takes a CreateIdentityProviderDetails interface as input
	rCreate := identity.CreateIdentityProviderRequest{}

	// The CreateSaml2IdentityProviderDetails struct implements the CreateIdentityProviderDetails interface
	details := identity.CreateSaml2IdentityProviderDetails{}
	details.CompartmentId = common.String(getTenancyID())
	details.Name = common.String("someName")
	//... more setup if needed
	// Use the above struct
	rCreate.CreateIdentityProviderDetails = details

	// Make the call
	rspCreate, createErr := c.CreateIdentityProvider(context.Background(), rCreate)

In the case of a polymorphic response you can type assert the interface to the expected type. For example:

	rRead := identity.GetIdentityProviderRequest{}
	rRead.IdentityProviderId = common.String("aValidId")
	response, err := c.GetIdentityProvider(context.Background(), rRead)

	provider := response.IdentityProvider.(identity.Saml2IdentityProvider)

An example of polymorphic json request handling can be found here: https://github.com/oracle/oci-go-sdk/blob/master/example/example_core_test.go#L63

# Pagination

When calling a list operation, the operation will retrieve a page of results. To retrieve more data, call the list operation again,
passing in the value of the most recent response's OpcNextPage as the value of Page in the next list operation call.
When there is no more data the OpcNextPage field will be nil. An example of pagination using this logic can be found here: https://github.com/oracle/oci-go-sdk/blob/master/example/example_core_pagination_test.go

# Logging and Debugging

The SDK has a built-in logging mechanism used internally. The internal logging logic is used to record the raw http
requests, responses and potential errors when (un)marshalling request and responses.

Built-in logging in the SDK is controlled via the environment variable "OCI_GO_SDK_DEBUG" and its contents. The below are possible values for the "OCI_GO_SDK_DEBUG" variable

1. "info" or "i" enables all info logging messages

2. "debug" or "d"  enables all debug and info logging messages

3. "verbose" or "v" or "1" enables all verbose, debug and info logging messages

4. "null" turns all logging messages off.

If the value of the environment variable does not match any of the above then default logging level is "info".
If the environment variable is not present then no logging messages are emitted.

The default destination for logging is Stderr and if you want to output log to a file you can set via environment variable "OCI_GO_SDK_LOG_OUTPUT_MODE". The below are possible values

1. "file" or "f" enables all logging output saved to file

2. "combine" or "c" enables all logging output to both stderr and file

# If the value does not match any of the above or does not exist then default logging output will be set to Stderr

You can also customize the log file location and name via "OCI_GO_SDK_LOG_FILE" environment variable, the value should be the path to a specific file
If this environment variable is not present, the default location will be the project root path

# Retry

Sometimes you may need to wait until an attribute of a resource, such as an instance or a VCN, reaches a certain state.
An example of this would be launching an instance and then waiting for the instance to become available, or waiting until a subnet in a VCN has been terminated.
You might also want to retry the same operation again if there's network issue etc...
This can be accomplished by using the RequestMetadata.RetryPolicy(request level configuration), alternatively, global(all services) or client level RetryPolicy configration is also possible.
You can find the examples here: https://github.com/oracle/oci-go-sdk/blob/master/example/example_retry_test.go

If you are trying to make a PUT/POST API call with binary request body, please make sure the binary request body is resettable, which means the request body should inherit Seeker interface.

The Retry behavior Precedence (Highest to lowest) is defined as below:-

	Operation level retry policy
	Client level retry policy
	Global level retry policy
	Environment level default retry policy for default retry
	Service level default retry policy

# Default Retry Policy

The OCI Go SDK defines a default retry policy that retries on the errors suitable for retries (see https://docs.oracle.com/en-us/iaas/Content/API/References/apierrors.htm),
for a recommended period of time (up to 7 attempts spread out over at most approximately 1.5 minutes). This default retry policy can be created using:

	// use SDK's default retry policy
	defaultRetryPolicy := common.DefaultRetryPolicy()

You can set this retry policy for a single request:

	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: &defaultRetryPolicy,
	}

or for all requests made by a client:

	client.SetCustomClientConfiguration(common.CustomClientConfiguration{
		RetryPolicy: &defaultRetryPolicy,
	})

or for all requests made by all clients:

	common.GlobalRetry = &defaultRetryPolicy

or setting default retry via environment varaible, which is a global switch for all services:

	export OCI_SDK_DEFAULT_RETRY_ENABLED=TRUE

Some services enable retry for operations by default, this can be overridden using any alternatives mentioned above.

# Eventual Consistency

Some resources may have to be replicated across regions and are only eventually consistent. That means the request to create, update, or delete the resource succeeded,
but the resource is not available everywhere immediately. Creating, updating, or deleting any resource in the Identity service is affected by eventual consistency, and
doing so may cause other operations in other services to fail until the Identity resource has been replicated.

For example, the request to CreateTag in the Identity service in the home region succeeds, but immediately using that created tag in another region in a request to
LaunchInstance in the Compute service may fail.

If you are creating, updating, or deleting resources in the Identity service, we recommend using an eventually consistent retry policy for any service you access. The
default retry policy already deals with eventual consistency. Example:

	// use SDK's default retry policy (which also deals with eventual consistency)
	defaultRetryPolicy := common.DefaultRetryPolicy()

This retry policy will use a different strategy if an eventually consistent change was made in the recent past (called the "eventually consistent window", currently
defined to be 4 minutes after the eventually consistent change). This special retry policy for eventual consistency will:

1. make up to 9 attempts (including the initial attempt); if an attempt is successful, no more attempts will be made

2. retry at most until (a) approximately the end of the eventually consistent window or (b) the end of the default retry period of about 1.5 minutes, whichever is
farther in the future; if an attempt is successful, no more attempts will be made, and the OCI Go SDK will not wait any longer

3. retry on the error codes 400-RelatedResourceNotAuthorizedOrNotFound, 404-NotAuthorizedOrNotFound, and 409-NotAuthorizedOrResourceAlreadyExists, for which the
default retry policy does not retry, in addition to the errors the default retry policy retries on (see https://docs.oracle.com/en-us/iaas/Content/API/References/apierrors.htm)

If there were no eventually consistent actions within the recent past, then this special retry strategy is not used.

If you want a retry policy that does not handle eventual consistency in a special way, for example because you retry on all error responses, you can use
DefaultRetryPolicyWithoutEventualConsistency or NewRetryPolicyWithOptions with the common.ReplaceWithValuesFromRetryPolicy(common.DefaultRetryPolicyWithoutEventualConsistency()) option:

	// use SDK's default retry policy, but without eventual consistency
	noEcRetryPolicy := common.DefaultRetryPolicyWithoutEventualConsistency()

	// or
	noEcRetryPolicy := common.NewRetryPolicyWithOptions(
		common.ReplaceWithValuesFromRetryPolicy(common.DefaultRetryPolicyWithoutEventualConsistency()),
		// possibly other options...
	)

The NewRetryPolicy function also creates a retry policy without eventual consistency.

# Circuit Breaker

Circuit Breaker can prevent an application repeatedly trying to execute an operation that is likely to fail, allowing it to continue without waiting for the fault to be rectified or wasting CPU cycles,
of course, it also enables an application to detect whether the fault has been resolved. If the problem appears to have been rectified, the application can attempt to invoke the operation.
Go SDK intergrates sony/gobreaker solution, wraps in a circuit breaker object, which monitors for failures. Once the failures reach a certain threshold, the circuit breaker trips,
and all further calls to the circuit breaker return with an error, this also saves the service from being overwhelmed with network calls in case of an outage.

Go SDK enable circuit breaker with default configuration, if you don't want to enable the solution, can disable the functionality before your application running
Go SDK also supports customize Circuit Breaker with specified configuratoins. You can find the examples here: https://github.com/oracle/oci-go-sdk/blob/master/example/example_circuitbreaker_test.go

# Using the SDK with a Proxy Server

The GO SDK uses the net/http package to make calls to OCI services. If your environment requires you to use a proxy server for outgoing HTTP requests
then you can set this up in the following ways:

1. Configuring environment variable as described here https://golang.org/pkg/net/http/#ProxyFromEnvironment
2. Modifying the underlying Transport struct for a service client

In order to modify the underlying Transport struct in HttpClient, you can do something similar to (sample code for audit service client):

	// create audit service client
	client, clerr := audit.NewAuditClientWithConfigurationProvider(common.DefaultConfigProvider())

	// create a proxy url
	proxyURL, err := url.Parse("http(s)://[username]:[password]@[ip address]:[port]")

	client.HTTPClient = &http.Client{
		// adding the proxy settings to the http.Transport
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

# Forward Compatibility

Some response fields are enum-typed. In the future, individual services may return values not covered by existing enums
for that field. To address this possibility, every enum-type response field is a modeled as a type that supports any string.
Thus if a service returns a value that is not recognized by your version of the SDK, then the response field will be set to this value.

When individual services return a polymorphic json response not available as a concrete struct, the SDK will return an implementation that only satisfies
the interface modeling the polymorphic json response.

# Contributions

Got a fix for a bug, or a new feature you'd like to contribute? The SDK is open source and accepting pull requests on GitHub
https://github.com/oracle/oci-go-sdk

# License

Licensing information available at: https://github.com/oracle/oci-go-sdk/blob/master/LICENSE.txt

# Notifications

To be notified when a new version of the Go SDK is released, subscribe to the following feed: https://github.com/oracle/oci-go-sdk/releases.atom

# Questions or Feedback

Please refer to this link: https://github.com/oracle/oci-go-sdk#help
*/

package oci

//go:generate go run cmd/genver/main.go cmd/genver/version_template.go --output common/version.go
