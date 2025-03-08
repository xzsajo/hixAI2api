package main

import (
	"fmt"
	"hixai2api/common"
)

func main() {
	text := "user_group=146; first-visit-url=https%3A%2F%2Fhix.ai%2Fzh; device-id=3ab3c15d54dda3a82aef9461f306e6fb; _gcl_au=1.1.1766983635.1741332968; _fbp=fb.1.1741332968166.96038000529332917; _tt_enable_cookie=1; _ttp=01JNQRGWQPYNM40F73N6A7DNWX_.tt.1; _ga=GA1.1.83506002.1741332969; FPID=FPID2.2.cveO53H73pFObELS1NXF0joKa3smasV1Fd7c6cbjXEs%3D.1741332969; g_state={\"i_l\":0}; __gsas=ID=0cca1082a003af27:T=1741333199:RT=1741333199:S=ALNI_MaEH-B7EYhLTgL2xoSVgNf0IhQJsw; __Secure-next-auth.callback-url=https%3A%2F%2Fhix.ai; __Host-next-auth.csrf-token=e4bbb2c59d9e74642b5d7926d5c479b82bb3d59319d93a703456e4e671769c5e%7C436ace319f2234970a324ea6255857f09db7a4ddacbf7891125959289e759c4b; _clck=1sjmshv%7C2%7Cfu1%7C0%7C1892; FPLC=XAGCIV0FsibiEV4aEmP8bf%2Blxcoj1o9OgcYVt1O2QGqayBBYTPM%2BjGHLoD%2FnnWcujJyAEpbGO44%2F8Q8DXKyUSfL5Nok7jVaIbkzAsb5tulYC%2FqkkIQUkMFHFx4JvBw%3D%3D; cf_clearance=LKNFRKszSzw2NLEwLyme_oNxmhclQZgsJ.rhJDJSP.c-1741420622-1.2.1.1-hCDm_uME5Vwjm_ggq9olNlgHo5EHWO7lFpaXq8gaRqB00YcABDG.rajHf35OsSUFxA62bEqrWva1rS8mRqnc6TuH1RwEatw38nyoyjMpxAZjgfNuGW4g7F4tcyQt.k4l4ndLNy3uJCaIsp3TdSHICVSrnHD382ybVRzNxvsMu454dTidi7FdIPSVe4FIpowHhji9JZ20vo4Mm4jpY7GsxHGw0bNNUelpZQX2gLV7lekxp93X5VqMpOj9xn_DVejK6azmoQtzTaYDOv1Rl_TvM6FbnKsjMYacmr6QJ2cCCNfnL5bVQFmiswmimyJoIxxw4NcCsxDCWCkt2g.MehCNAqLrJEM82s7cmmfhhaqwwoKnfWnCofiFOgmWdHISWrXVxfkDP5f_lY4eaaZu7HShbeKoNLTGBQAHGoX1ZzklpEY; user-info=%7B%22id%22%3A%22cm7ygtgt9005rczp0k9klmjo2%22%2C%22name%22%3A%22Dean%20Xv%22%2C%22email%22%3A%22alanxv1024%40gmail.com%22%2C%22firstName%22%3A%22Dean%22%2C%22lastName%22%3A%22Xv%22%2C%22image%22%3A%22https%3A%2F%2Flh3.googleusercontent.com%2Fa%2FACg8ocJho_JMw-_dArGBTNa24ipeEGZ8dzv9ttGFTO1A7VWlc7p7Gw%3Ds96-c%22%7D; ph_phc_9LvbXawaTFdrUSVPTOjwbVv7bZWE1iOQDhF8U7dPa0E_posthog=%7B%22distinct_id%22%3A%22cm7ygtgt9005rczp0k9klmjo2%22%2C%22%24sesid%22%3A%5B1741420623114%2C%22019574c0-a322-730d-b08a-5e060a804374%22%2C1741420536610%5D%7D; _uetsid=d5317470fb2611efa28ab1f1448f428f; _uetvid=d531a870fb2611ef8e2ccd156937b0a5; _ga_NSMTS5QKXE=GS1.1.1741426506.4.1.1741426506.0.0.0; _ga_JS0YXNKF22=GS1.1.1741426506.4.1.1741426506.60.0.0; __Secure-next-auth.session-token=eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIn0..bHsl69s76Gt9pDlU.Z0r4efSQQC7_ApzfG52auu3zxwITxxW2oSsG_koOu5LdWabLijYTNerBvv9igdNi3UOg6pWsN0WLW6ES0CbSskNQ44_c2D8dkPOZk5TsMTnOz-kiHRqnQc1NXkQX-pbcjQNYVQY8MZ1EKce6gMdx4MBtJNXKuEKhRV0DgD-4RAQlYsVziHklVYhbV78b1aG1EHpkxEcdWDea3ORGVvXAct4K_qtovuudmzaPAlg68IMC3jO1vnLRFveAW-12AYm_KbSvBtuhz7zipwL_KUC5NRB7dD1XB3fb2m3z7hFZA3epPulH1Ub2sZZC_QCwRYpnEf4UqkX9MxiTpJ0V-Gz2mGGiyN5tJ8Y4LgjC4DDo8griXpgyMkW-5tCsgQFAZaZLwsy93PbswoWaOotGLIEbHkzxJ87t-RviJZIROF569NP0I_PAz80Floo2-iOfGOswbWh9uPKbte7mJ8Y1q7Kmz5G4hMyTYfROl_L-NUC9p13RlJYvN1leCVDJpIK6ttVrmTX9-rr6nbs.PSmprRO2wnLYO1OFTbgWZA"
	//text := "[{\"role\":\"user\",\"content\":\"hi\"},{\"role\":\"assistant\",\"content\":\"Hello! How can I assist you today?\"}]"

	sha256Hash := common.StringToSHA256(text)
	//dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f
	//dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f
	//fmt.Printf("MD5: %s\n", md5Hash)
	//fmt.Printf("SHA1: %s\n", sha1Hash)
	fmt.Printf("SHA256: %s\n", sha256Hash)
}
