//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package connection

import (
	"github.com/flanksource/duty/types"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSConnection) DeepCopyInto(out *AWSConnection) {
	*out = *in
	in.AccessKey.DeepCopyInto(&out.AccessKey)
	in.SecretKey.DeepCopyInto(&out.SecretKey)
	in.SessionToken.DeepCopyInto(&out.SessionToken)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSConnection.
func (in *AWSConnection) DeepCopy() *AWSConnection {
	if in == nil {
		return nil
	}
	out := new(AWSConnection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Authentication) DeepCopyInto(out *Authentication) {
	*out = *in
	in.Username.DeepCopyInto(&out.Username)
	in.Password.DeepCopyInto(&out.Password)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Authentication.
func (in *Authentication) DeepCopy() *Authentication {
	if in == nil {
		return nil
	}
	out := new(Authentication)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GCPConnection) DeepCopyInto(out *GCPConnection) {
	*out = *in
	if in.Credentials != nil {
		in, out := &in.Credentials, &out.Credentials
		*out = new(types.EnvVar)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GCPConnection.
func (in *GCPConnection) DeepCopy() *GCPConnection {
	if in == nil {
		return nil
	}
	out := new(GCPConnection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GCSConnection) DeepCopyInto(out *GCSConnection) {
	*out = *in
	in.GCPConnection.DeepCopyInto(&out.GCPConnection)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GCSConnection.
func (in *GCSConnection) DeepCopy() *GCSConnection {
	if in == nil {
		return nil
	}
	out := new(GCSConnection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *S3Connection) DeepCopyInto(out *S3Connection) {
	*out = *in
	in.AWSConnection.DeepCopyInto(&out.AWSConnection)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new S3Connection.
func (in *S3Connection) DeepCopy() *S3Connection {
	if in == nil {
		return nil
	}
	out := new(S3Connection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SFTPConnection) DeepCopyInto(out *SFTPConnection) {
	*out = *in
	in.Authentication.DeepCopyInto(&out.Authentication)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SFTPConnection.
func (in *SFTPConnection) DeepCopy() *SFTPConnection {
	if in == nil {
		return nil
	}
	out := new(SFTPConnection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SMBConnection) DeepCopyInto(out *SMBConnection) {
	*out = *in
	in.Authentication.DeepCopyInto(&out.Authentication)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SMBConnection.
func (in *SMBConnection) DeepCopy() *SMBConnection {
	if in == nil {
		return nil
	}
	out := new(SMBConnection)
	in.DeepCopyInto(out)
	return out
}
