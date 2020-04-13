package adpcm

// -------------------------------------------------------------------------------------

var indexAdjust = [8]int{
	-1, -1, -1, -1, 2, 4, 6, 8,
}

var stepTable = [89]int{
	7, 8, 9, 10, 11, 12, 13, 14, 16, 17,
	19, 21, 23, 25, 28, 31, 34, 37, 41, 45,
	50, 55, 60, 66, 73, 80, 88, 97, 107, 118,
	130, 143, 157, 173, 190, 209, 230, 253, 279, 307,
	337, 371, 408, 449, 494, 544, 598, 658, 724, 796,
	876, 963, 1060, 1166, 1282, 1411, 1552, 1707, 1878, 2066,
	2272, 2499, 2749, 3024, 3327, 3660, 4026, 4428, 4871, 5358,
	5894, 6484, 7132, 7845, 8630, 9493, 10442, 11487, 12635, 13899,
	15289, 16818, 18500, 20350, 22385, 24623, 27086, 29794, 32767,
}

type status struct {
	sample int
	index  int
}

func (s *status) decodeSample(scode byte) int {
	// 将 scode 分离为数据和符号
	code := int(scode & 7)
	delta := ((stepTable[s.index] * code) >> 2) + (stepTable[s.index] >> 3) // 后面加的一项是为了减少误差
	if (scode & 8) != 0 {
		delta = -delta // 负数
	}
	s.sample += delta // 计算出当前的波形数据
	s.index += indexAdjust[code]
	if s.index < 0 {
		s.index = 0
	} else if s.index > 88 {
		s.index = 88
	}
	if s.sample > 32767 {
		return 32767
	}
	if s.sample < -32768 {
		return -32768
	}
	return s.sample
}

func (s *status) saveSample(samples []byte, idx int, scode byte) int {
	sample := s.decodeSample(scode)
	return saveSample(samples, idx, sample)
}

func saveSample(samples []byte, idx int, sample int) int {
	samples[idx] = byte(sample)
	samples[idx+1] = byte(sample >> 8)
	return idx + 2
}

// -------------------------------------------------------------------------------------

func loadStatus(b []byte) *status {
	return &status{
		sample: int(int16(b[0]) | (int16(b[1]) << 8)),
		index:  int(b[2]),
	}
}

func loadBlock(channels int, block []byte, samples []byte) {
	status1 := loadStatus(block)
	status2 := status1
	saveSample(samples, 0, status1.sample)
	if channels > 1 {
		status2 = loadStatus(block[4:])
		saveSample(samples, 2, status2.sample)
	}
	idx := channels << 1
	for _, code := range block[(channels << 2):] {
		idx = status1.saveSample(samples, idx, code>>4)
		idx = status2.saveSample(samples, idx, code&0xf)
	}
}

// -------------------------------------------------------------------------------------
