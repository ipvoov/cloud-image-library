package picture

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/consts"
	"cloud/internal/model/entity"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/lucasb-eyer/go-colorful"
)

// parseImageInfo 解析图片信息
func (s *sPicture) parseImageInfo(ctx context.Context, file *ghttp.UploadFile) (width, height int, format string, scale float64, err error) {
	// 打开文件
	f, err := file.Open()
	if err != nil {
		g.Log().Errorf(ctx, "打开文件失败: %v", err)
		return 0, 0, "", 0, err
	}
	defer f.Close()

	// 解析图片配置信息
	config, format, err := image.DecodeConfig(f)
	if err != nil {
		g.Log().Errorf(ctx, "解析图片配置失败: %v", err)
		// 如果解析失败，尝试从文件扩展名获取格式
		format = s.getFormatFromFilename(file.Filename)
		return 0, 0, format, 0, err
	}

	width = config.Width
	height = config.Height

	// 计算宽高比例
	if height > 0 {
		scale = float64(width) / float64(height)
	}

	g.Log().Infof(ctx, "解析图片信息成功: %dx%d, 格式: %s, 比例: %.2f", width, height, format, scale)

	return width, height, format, scale, nil
}

// getFormatFromFilename 从文件名获取格式
func (s *sPicture) getFormatFromFilename(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "jpeg"
	case ".png":
		return "png"
	case ".gif":
		return "gif"
	case ".bmp":
		return "bmp"
	case ".webp":
		return "webp"
	default:
		return "jpeg" // 默认格式
	}
}

// parseImageInfoFromURL 从URL解析图片信息
func (s *sPicture) parseImageInfoFromURL(ctx context.Context, url string) (width, height int, format string, scale float64, fileSize int64, err error) {
	// 先从URL路径获取格式作为备用
	format = s.getFormatFromFilename(url)

	g.Log().Infof(ctx, "开始解析URL图片信息: %s", url)

	// 创建HTTP客户端，设置超时时间
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// 发送HTTP请求获取图片
	resp, err := client.Get(url)
	if err != nil {
		g.Log().Errorf(ctx, "获取URL图片失败: %v", err)
		return 0, 0, format, 0, 0, err
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("HTTP状态码错误: %d", resp.StatusCode)
		g.Log().Errorf(ctx, "HTTP请求失败: %v", err)
		return 0, 0, format, 0, 0, err
	}

	// 获取文件大小
	fileSize = resp.ContentLength
	if fileSize < 0 {
		// 如果Content-Length头不存在或为负数，尝试从响应体计算大小
		// 注意：这会消耗响应体，需要重新请求来获取图片数据
		g.Log().Warningf(ctx, "无法从Content-Length获取文件大小，将重新请求")
		fileSize = 0 // 暂时设为0，后续可以通过其他方式获取
	}

	// 解析图片配置信息
	config, detectedFormat, err := image.DecodeConfig(resp.Body)
	if err != nil {
		g.Log().Errorf(ctx, "解析图片配置失败: %v", err)
		// 如果解析失败，返回从文件名获取的格式
		return 0, 0, format, 0, fileSize, err
	}

	width = config.Width
	height = config.Height
	format = detectedFormat

	// 计算宽高比例
	if height > 0 {
		scale = float64(width) / float64(height)
	}

	g.Log().Infof(ctx, "成功解析URL图片信息: %dx%d, 格式: %s, 比例: %.2f, 大小: %d bytes", width, height, format, scale, fileSize)
	return width, height, format, scale, fileSize, nil
}

// entityToVO 将entity转换为VO
func (s *sPicture) entityToVO(ctx context.Context, picture *entity.Picture) *v1.PictureVO {
	// 解析标签JSON
	var tags []string
	if picture.Tags != "" {
		err := gjson.DecodeTo(picture.Tags, &tags)
		if err != nil {
			// 如果JSON解析失败，尝试按逗号分割
			tags = strings.Split(picture.Tags, ",")
			for i, tag := range tags {
				tags[i] = strings.TrimSpace(tag)
			}
		}
	}

	return &v1.PictureVO{
		Id:           picture.Id,
		Url:          picture.Url,
		Name:         picture.Name,
		Introduction: picture.Introduction,
		Category:     picture.Category,
		Tags:         tags,
		PicSize:      picture.PicSize,
		PicWidth:     picture.PicWidth,
		PicHeight:    picture.PicHeight,
		PicScale:     picture.PicScale,
		PicFormat:    picture.PicFormat,
		UserId:       picture.UserId,
		SpaceId:      picture.SpaceId,
		CreateTime:   picture.CreateTime.Format(consts.Y_m_d_His),
		EditTime:     picture.EditTime.Format(consts.Y_m_d_His),
		UpdateTime:   picture.UpdateTime.Format(consts.Y_m_d_His),
		ThumbnailUrl: picture.ThumbnailUrl,
		PicColor:     picture.PicColor,
	}
}

// entityToPicture 将entity转换为Picture（管理员视图）
func (s *sPicture) entityToPicture(ctx context.Context, picture *entity.Picture) *v1.Picture {
	// 解析标签JSON - 前端期望JSON字符串格式
	var tagsJson string
	if picture.Tags != "" {
		// 如果数据库中已经是JSON字符串，直接使用
		tagsJson = picture.Tags
	} else {
		// 如果为空，返回空数组的JSON字符串
		tagsJson = "[]"
	}

	// 处理时间格式
	var createTime, editTime, updateTime, reviewTime string
	if picture.CreateTime != nil {
		createTime = picture.CreateTime.Format("Y-m-d H:i:s")
	}
	if picture.EditTime != nil {
		editTime = picture.EditTime.Format("Y-m-d H:i:s")
	}
	if picture.UpdateTime != nil {
		updateTime = picture.UpdateTime.Format("Y-m-d H:i:s")
	}
	if picture.ReviewTime != nil {
		reviewTime = picture.ReviewTime.Format("Y-m-d H:i:s")
	}

	return &v1.Picture{
		Id:            picture.Id,
		Url:           picture.Url,
		Name:          picture.Name,
		Introduction:  picture.Introduction,
		Category:      picture.Category,
		Tags:          tagsJson, // 前端期望JSON字符串格式
		PicSize:       picture.PicSize,
		PicWidth:      picture.PicWidth,
		PicHeight:     picture.PicHeight,
		PicScale:      picture.PicScale,
		PicFormat:     picture.PicFormat,
		UserId:        picture.UserId,
		SpaceId:       picture.SpaceId,
		CreateTime:    createTime,
		EditTime:      editTime,
		UpdateTime:    updateTime,
		ThumbnailUrl:  picture.ThumbnailUrl,
		PicColor:      picture.PicColor,
		IsDelete:      picture.IsDelete,
		ReviewStatus:  picture.ReviewStatus,
		ReviewMessage: picture.ReviewMessage,
		ReviewerId:    picture.ReviewerId,
		ReviewTime:    reviewTime,
	}
}

// extractDominantColor 提取图片主色调
func (s *sPicture) extractDominantColor(ctx context.Context, file *ghttp.UploadFile) (string, error) {
	// 打开文件
	f, err := file.Open()
	if err != nil {
		g.Log().Errorf(ctx, "打开文件失败: %v", err)
		return "", err
	}
	defer f.Close()

	// 解码图片
	img, _, err := image.Decode(f)
	if err != nil {
		g.Log().Errorf(ctx, "解码图片失败: %v", err)
		return "", err
	}

	// 使用优化的主色调提取算法
	return s.extractDominantColorOptimized(img, ctx)
}

// extractDominantColorFromURL 从URL提取图片主色调
func (s *sPicture) extractDominantColorFromURL(ctx context.Context, url string) (string, error) {

	// 创建HTTP客户端，设置超时时间
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// 发送HTTP请求获取图片
	resp, err := client.Get(url)
	if err != nil {
		g.Log().Errorf(ctx, "获取URL图片失败: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("HTTP状态码错误: %d", resp.StatusCode)
		g.Log().Errorf(ctx, "HTTP请求失败: %v", err)
		return "", err
	}

	// 解码图片
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		g.Log().Errorf(ctx, "解码URL图片失败: %v", err)
		return "", err
	}

	// 使用优化的主色调提取算法
	return s.extractDominantColorOptimized(img, ctx)
}

// extractDominantColorOptimized 主色调提取核心算法（使用prominentcolor库）
func (s *sPicture) extractDominantColorOptimized(img image.Image, ctx context.Context) (string, error) {
	// 使用prominentcolor库提取主色调
	// 参数说明：
	// - 1: 提取1个主色调
	// - img: 图片对象
	// - 0: 不裁剪图片
	// - 0: 背景掩码
	// - nil: 无额外参数
	colors, err := prominentcolor.KmeansWithAll(1, img, 0, 0, nil)
	if err != nil {
		g.Log().Warningf(ctx, "prominentcolor提取主色调失败，使用备用算法: %v", err)
		return s.extractDominantColorFallback(img, ctx)
	}

	if len(colors) == 0 {
		g.Log().Warningf(ctx, "prominentcolor未提取到颜色，使用备用算法")
		return s.extractDominantColorFallback(img, ctx)
	}

	// 获取第一个（最主要的）颜色
	dominantColor := colors[0]

	// 转换为十六进制格式
	hexColor := fmt.Sprintf("#%02X%02X%02X", dominantColor.Color.R, dominantColor.Color.G, dominantColor.Color.B)

	// 使用go-colorful库验证颜色有效性
	if _, err := colorful.Hex(hexColor); err != nil {
		g.Log().Warningf(ctx, "提取的颜色无效，使用备用算法: %v", err)
		return s.extractDominantColorFallback(img, ctx)
	}

	g.Log().Infof(ctx, "成功提取主色调: %s (权重: %d)", hexColor, dominantColor.Cnt)
	return hexColor, nil
}

// extractDominantColorFallback 备用主色调提取算法
func (s *sPicture) extractDominantColorFallback(img image.Image, ctx context.Context) (string, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 计算采样步长（更智能的采样）
	step := 1
	if width > 300 || height > 300 {
		step = 6 // 大图片每6个像素采样一次
	} else if width > 100 || height > 100 {
		step = 3 // 中等图片每3个像素采样一次
	}

	// 使用更高效的颜色统计
	colorCounts := make(map[string]int)
	totalPixels := 0

	// 采样图片像素
	for y := bounds.Min.Y; y < bounds.Max.Y; y += step {
		for x := bounds.Min.X; x < bounds.Max.X; x += step {
			r, g, b, _ := img.At(x, y).RGBA()

			// 量化颜色到更少的级别（16种颜色级别）
			r8 := int(r>>8) & 0xF0 // 只保留高4位
			g8 := int(g>>8) & 0xF0
			b8 := int(b>>8) & 0xF0

			colorKey := fmt.Sprintf("%02X%02X%02X", r8, g8, b8)
			colorCounts[colorKey]++
			totalPixels++
		}
	}

	// 找出最频繁的颜色
	var dominantColor string
	maxCount := 0
	for color, count := range colorCounts {
		if count > maxCount {
			maxCount = count
			dominantColor = color
		}
	}

	if dominantColor == "" {
		return "#000000", nil
	}

	// 转换为标准十六进制格式
	hexColor := "#" + dominantColor

	// 使用go-colorful库验证颜色有效性
	if _, err := colorful.Hex(hexColor); err != nil {
		g.Log().Warningf(ctx, "备用算法提取的颜色无效，使用默认黑色: %v", err)
		return "#000000", nil
	}

	g.Log().Infof(ctx, "备用算法成功提取主色调: %s", hexColor)
	return hexColor, nil
}

// RGB 颜色结构
type RGB struct {
	R, G, B int
}

// ColorSimilarityItem 颜色相似度项目
type ColorSimilarityItem struct {
	Picture    *entity.Picture
	Similarity float64
}

// parseHexColor 解析十六进制颜色值
func (s *sPicture) parseHexColor(hexColor string) (RGB, error) {
	// 移除#号前缀
	hexColor = strings.TrimPrefix(hexColor, "#")

	// 确保是6位十六进制
	if len(hexColor) != 6 {
		return RGB{}, fmt.Errorf("颜色格式错误，需要6位十六进制值")
	}

	// 解析RGB值
	r, err := strconv.ParseInt(hexColor[0:2], 16, 64)
	if err != nil {
		return RGB{}, err
	}

	g, err := strconv.ParseInt(hexColor[2:4], 16, 64)
	if err != nil {
		return RGB{}, err
	}

	b, err := strconv.ParseInt(hexColor[4:6], 16, 64)
	if err != nil {
		return RGB{}, err
	}

	return RGB{R: int(r), G: int(g), B: int(b)}, nil
}

// calculateEuclideanDistance 计算两个RGB颜色的欧氏距离相似度
func (s *sPicture) calculateEuclideanDistance(color1, color2 RGB) float64 {
	// 计算RGB欧氏距离
	dr := float64(color1.R - color2.R)
	dg := float64(color1.G - color2.G)
	db := float64(color1.B - color2.B)

	// 计算欧氏距离
	distance := math.Sqrt(dr*dr + dg*dg + db*db)

	// 转换为相似度（距离越小，相似度越高）
	// 最大可能距离是 sqrt(255^2 + 255^2 + 255^2) ≈ 441.67
	maxDistance := math.Sqrt(255*255 + 255*255 + 255*255)
	similarity := 1.0 - (distance / maxDistance)

	// 确保相似度在0-1范围内
	if similarity < 0 {
		similarity = 0
	}
	if similarity > 1 {
		similarity = 1
	}

	return similarity
}

// calculateColorSimilarity 计算颜色相似度并排序
func (s *sPicture) calculateColorSimilarity(ctx context.Context, pictures []entity.Picture, targetRGB RGB) []ColorSimilarityItem {
	var items []ColorSimilarityItem

	for _, picture := range pictures {
		if picture.PicColor == "" {
			continue // 跳过没有主色调的图片
		}

		// 解析图片的主色调
		pictureRGB, err := s.parseHexColor(picture.PicColor)
		if err != nil {
			g.Log().Warningf(ctx, "解析图片主色调失败: %s, 错误: %v", picture.PicColor, err)
			continue
		}

		// 计算相似度
		similarity := s.calculateEuclideanDistance(targetRGB, pictureRGB)

		items = append(items, ColorSimilarityItem{
			Picture:    &picture,
			Similarity: similarity,
		})
	}

	// 按相似度降序排序（相似度高的在前）
	sort.Slice(items, func(i, j int) bool {
		return items[i].Similarity > items[j].Similarity
	})

	g.Log().Infof(ctx, "计算颜色相似度完成，共处理 %d 张图片", len(items))

	return items
}
