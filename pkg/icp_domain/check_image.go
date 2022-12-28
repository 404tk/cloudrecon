package icp_domain

/*
func (i *Icp) checkImage() (string, error) {
	// 获取验证图像,UUID
	result := &IcpResponse{Params: &ImageParams{}}
	err := i.post("image/getCheckImage", nil, "application/x-www-form-urlencoded;charset=UTF-8", i.token, result)
	if err != nil {
		return "",err
	}

	if !result.Success {
		return "",fmt.Errorf("获取验证图片失败：%s", result.Msg)
	}

	params := result.Params.(*ImageParams)
	imgSrc, err := toBase64Img(params.BigImage)
	if err != nil {
		return "",fmt.Errorf("解析背景图失败:%w", err)
	}
	defer func(imgSrc *gocv.Mat) {
		_ = imgSrc.Close()
	}(imgSrc)

	imgTempl, err := toBase64Img(params.SmallImage)
	if err != nil {
		return "",fmt.Errorf("解析模块图失败:%w", err)
	}
	defer func(imgTempl *gocv.Mat) {
		_ = imgTempl.Close()
	}(imgTempl)

	res := gocv.NewMat()
	defer func(result *gocv.Mat) {
		_ = res.Close()
	}(&res)

	m := gocv.NewMat()
	gocv.MatchTemplate(*imgTempl, *imgSrc, &res, gocv.TmCcoeffNormed, m)
	defer func(m *gocv.Mat) {
		_ = m.Close()
	}(&m)

	_, maxConfidence, _, maxLoc := gocv.MinMaxLoc(res)
	fmt.Println(maxConfidence)
	if maxConfidence < 0.8 {
		return "",fmt.Errorf("图片验证可靠性过低:%f", maxConfidence)
	}

	// 通过拼图验证，获取sign
	signresp := &IcpResponse{}
	imgbody := fmt.Sprintf(`{"key":"%s","value":"%d"}`, params.Uuid, maxLoc.X+1)
	err = i.post("image/checkImage", bytes.NewReader([]byte(imgbody)), "application/json", i.token, signresp)
	if err != nil {
		return "",fmt.Errorf("验证图片失败:%w", err)
	}
	sign := signresp.Params.(string)
	fmt.Println(sign)

	return sign,nil
}

func toBase64Img(in string) (*gocv.Mat, error) {
	ddd, _ := base64.StdEncoding.DecodeString(in)
	imgSrc, err := gocv.IMDecode(ddd, gocv.IMReadGrayScale)
	if err != nil {
		return nil, err
	}
	if imgSrc.Empty() {
		return nil, errors.New("invalid read")
	}
	return &imgSrc, nil
}
*/
