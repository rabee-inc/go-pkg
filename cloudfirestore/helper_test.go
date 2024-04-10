package cloudfirestore_test

import (
	"testing"

	"github.com/rabee-inc/go-pkg/cloudfirestore"
)

func Test_ValidateDocumentID(t *testing.T) {
	type args struct {
		id string
	}
	type want struct {
		isValid bool
	}
	type testCase struct {
		name string
		args args
		want want
	}

	// テストケースの定義
	tcs := []testCase{
		{
			name: "正常系",
			args: args{
				id: "jXfPGTXBk1IWlITJ6cKIaDSVs8l8xY5cp5rS1Apzg0hgr4qK4PjYP3FjBusKB4W8wBvCG5JfngYvQ1SCDoO4t6EMZPgcPGUe7m4d",
			},
			want: want{
				isValid: true,
			},
		},
		{
			name: "異常系: 空",
			args: args{
				id: "",
			},
			want: want{
				isValid: false,
			},
		},
		{
			name: "異常系: 1,500 バイト以下にする必要があります",
			args: args{
				id: "LzYmgJPYwU8OLl2jRhjIpeVKux1bExrS6ezImLB7DzXqLg0tjDFhtMdtbJK2m61rtOLIMgDBrlvoJ58TIEXZor539TPDwINeIGpAEED8PKVTZ56p3GbPCkE28mVJVPytetebqmPOhia4trrjmTp34rlGAwcDsgrTSUSV7oaE2jjLgWDM9vBWV4oA5H5hCYEHGP57jt8RZnh5sjZKh8Llfd9R8DulGFCOoSsryReSqcZnnq3M0BnWKF0pRQsZiwTAKE2JepK8LZPVhNJ5mEZM1Yterl52WXVy1PILOMN4WTWOSR51VD2ycjSLgLdTDjEhGofrKOYEII8I3h8IUQzAjyFQlDZaliMzGxNt61gP7hShF30IfAtEUmJKt56r5kJ9q0qnUUsGZ0CMNN0oCkG24QXvU9AHQmRd4y6bPxsI3lsI71gkEA5dYfKOmK9qz0MCCfKKfchrARqPQzJ6uBGk7zvr7q4uM2SgpBIWKMO1kZVFVSELlJz1B6ztw0mvAjQpdRpTs4a0zPAYBS43xyiNBsOTXAUIh1BFnTzN3xQ2I4cBx1gQTanwfxfvCpj2aB1434yq6AQoVtqgCvXSd0SspBZ1RDwGWdv5zotJUPffN4Nww5HmCYQUMJ04TEVVz5ySU1Ri0SeVgoDugWNiB0i74ftp2xgvTwYLNMRPIq3QFh7euLVJLJVmrbtb9iD5pfjmSVN2hb20IiU0cVFzSiRghxx7BixmUPu5CYtGhICAgf0C6qm68eKfsXpQCP7u3d1vXVgQgn8Wlvf7JjwRUEagnXnH8tnfJjA4tIENygH71r7XBiiUnNtbPlCHGcLqdWE5nImBFZtWpvoe8wFOlfSkOkmiMUwkZQVCQh49Hd8X5rAtsbSSYWTGdAI2rFOZpNe9tap8no9GjGMsbJ2rnzAzI1u7atJZhgriKdT8CakexmU7KOIkVMHVELSxTunjemewJx10Q8DMYRi36KmN5BQy5RAUEwTB1dfZHOnVLi2Oe3MdFGRfaMdIF3ogC253T7ruQniaFPf6QjLhS2ZFSI9INwjjYpP19tpsq3Y15obt94zPJNZfvtzuTErWodFFHRSHTLiTxBoAVD3Q4W4WlBHJO0eDpPt2HNyaneD2SWkSc6OqVYl0hEI2PVSV0iGflGZNjdNcrM9BJLaGB3gRJbBbgaiqtkaUnyacD9WxmEQjGCMJqKmONRz6Q9vsA4h2pPdjslmTchhh4zir13tQsJL7JOa0XkcUkZF9f4rlptUVAV4Eniro0iyW8KZcc6elOekRx5cCFq6K9cYPG3DIiP9xrGykWoa7qcv3GdcOdZRYwD1VhQuGqzMdZD7N2y6fYQ7Lcby3nf8tix2jPXjvbqsiExKo4ZHxjhYfldBY80lFfXjNDFvcxRsBrRdHe5y0kM5iXDhW36rsX2UNqgU06WXm8zvYenwmGeKl7oiSFNbmw9VVg0zC7JYrRK4RZv3CoFo9cuZ0Sx6hGZvTZAklculZhvpen1I06",
			},
			want: want{
				isValid: false,
			},
		},
		{
			name: "異常系: スラッシュ（/）は使用できません",
			args: args{
				id: "hoge/hoge",
			},
			want: want{
				isValid: false,
			},
		},
		{
			name: "異常系: 1 つのピリオド（.）または 2 つのピリオド（..）のみで構成することはできません。",
			args: args{
				id: "..hoge",
			},
			want: want{
				isValid: false,
			},
		},
		{
			name: "異常系: 次の正規表現とは照合できません: __.*__",
			args: args{
				id: "__hoge__",
			},
			want: want{
				isValid: false,
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			isValid := cloudfirestore.ValidateDocumentID(tc.args.id)
			if isValid != tc.want.isValid {
				t.Errorf("got: %v, want: %v", isValid, tc.want.isValid)
			}
		})
	}
}
