package svgutil_test

import (
	"testing"

	"github.com/leclerc04/go-tool/agl/util/svgutil"
	"github.com/stretchr/testify/assert"
)

var inTheEnd = `In The End -- Linkin Park(林肯公园)
It starts with one thing
I don't know why
It doesn't even matter how hard you try

I tried so hard
And got so far
But in the end
It doesn't even matter
I had to fall
To lose it all
But in the end
It doesn't even matter`

func Test(t *testing.T) {
	text2pathM := svgutil.NewManager()
	ret := text2pathM.TextToPath(inTheEnd, 0, 20)

	{
		assert.Contains(t, ret, `<svg height="260.000000px" width="380.888672px">
<g transform="scale(0.009766) translate(0,2048.000000)">
<path transform="translate(0.000000) rotate(180) scale(-1, 1)" d="M549 1380h-149v-1211h149v-169h-498v169h149v1211h-149v169h498v-169z" />
<path transform="translate(602.000000) rotate(180) scale(-1, 1)" d="M1110 0h-195v629q0 341 -249 341q-129 0 -212 -96t-83 -243v-631h-196v1106h196v-183h4q125 209 363 209q182 0 277 -116.5t95 -339.5v-676z" />
<path transform="translate(1864.000000) rotate(180) scale(-1, 1)" d="" />
<path transform="translate(2470.000000) rotate(180) scale(-1, 1)" d="M1136 1371h-447v-1371h-201v1371h-445v178h1093v-178z" />
<path transform="translate(3644.000000) rotate(180) scale(-1, 1)" d="M1109 0h-195v636q0 334 -248 334q-127 0 -211 -97.5t-84 -245.5v-627h-196v1637h196v-713h4q127 208 361 208q373 0 373 -451v-681z" />
<path transform="translate(4905.000000) rotate(180) scale(-1, 1)" d="M1076 503h-774q5 -177 97 -273.5t257 -96.5q186 0 341 119v-175q-146 -103 -387 -103q-240 0 -374.5 152t-134.5 422q0 254 147.5 419t366.5 165q217 0 339 -139.5t122 -390.5v-99zM880 660q-1 150 -71 232t-197 82q-117 0 -202.5 -87t-107.5 -227h578z" />
<path transform="translate(6067.000000) rotate(180) scale(-1, 1)" d="" />
<path transform="translate(6673.000000) rotate(180) scale(-1, 1)" d="M1033 0h-833v1549h798v-178h-598v-496h554v-177h-554v-521h633v-177z" />
<path transform="translate(7799.000000) rotate(180) scale(-1, 1)" d="M1110 0h-195v629q0 341 -249 341q-129 0 -212 -96t-83 -243v-631h-196v1106h196v-183h4q125 209 363 209q182 0 277 -116.5t95 -339.5v-676z" />
<path transform="translate(9061.000000) rotate(180) scale(-1, 1)" d="M1135 0h-196v184h-4q-121 -210 -376 -210q-210 0 -334 149.5t-124 403.5q0 274 140 439.5t366 165.5t328 -176h4v681h196v-1637zM940 658q0 133 -87.5 222.5t-215.5 89.5q-156 0 -246 -116.5t-90 -317.5q0 -185 86.5 -292.5t230.5 -107.5q142 0 232 105.5t90 262.5v154z
" />`)
	}

	ret = text2pathM.TextToPath(inTheEnd, 100, 20)
	{
		assert.Contains(t, ret, `<svg height="580.000000px" width="100.000000px">
<g transform="scale(0.009766) translate(0,2048.000000)">
<path transform="translate(0.000000) rotate(180) scale(-1, 1)" d="M549 1380h-149v-1211h149v-169h-498v169h149v1211h-149v169h498v-169z" />
<path transform="translate(602.000000) rotate(180) scale(-1, 1)" d="M1110 0h-195v629q0 341 -249 341q-129 0 -212 -96t-83 -243v-631h-196v1106h196v-183h4q125 209 363 209q182 0 277 -116.5t95 -339.5v-676z" />
<path transform="translate(1864.000000) rotate(180) scale(-1, 1)" d="" />
<path transform="translate(2470.000000) rotate(180) scale(-1, 1)" d="M1136 1371h-447v-1371h-201v1371h-445v178h1093v-178z" />
<path transform="translate(3644.000000) rotate(180) scale(-1, 1)" d="M1109 0h-195v636q0 334 -248 334q-127 0 -211 -97.5t-84 -245.5v-627h-196v1637h196v-713h4q127 208 361 208q373 0 373 -451v-681z" />
<path transform="translate(4905.000000) rotate(180) scale(-1, 1)" d="M1076 503h-774q5 -177 97 -273.5t257 -96.5q186 0 341 119v-175q-146 -103 -387 -103q-240 0 -374.5 152t-134.5 422q0 254 147.5 419t366.5 165q217 0 339 -139.5t122 -390.5v-99zM880 660q-1 150 -71 232t-197 82q-117 0 -202.5 -87t-107.5 -227h578z" />
<path transform="translate(6067.000000) rotate(180) scale(-1, 1)" d="" />
<path transform="translate(6673.000000) rotate(180) scale(-1, 1)" d="M1033 0h-833v1549h798v-178h-598v-496h554v-177h-554v-521h633v-177z" />
<path transform="translate(7799.000000) rotate(180) scale(-1, 1)" d="M1110 0h-195v629q0 341 -249 341q-129 0 -212 -96t-83 -243v-631h-196v1106h196v-183h4q125 209 363 209q182 0 277 -116.5t95 -339.5v-676z" />
</g>
<g transform="scale(0.009766) translate(0,4096.000000)">
<path transform="translate(0.000000) rotate(180) scale(-1, 1)" d="M1135 0h-196v184h-4q-121 -210 -376 -210q-210 0 -334 149.5t-124 403.5q0 274 140 439.5t366 165.5t328 -176h4v681h196v-1637zM940 658q0 133 -87.5 222.5t-215.5 89.5q-156 0 -246 -116.5t-90 -317.5q0 -185 86.5 -292.5t230.5 -107.5q142 0 232 105.5t90 262.5v154z
" />
<path transform="translate(1310.000000) rotate(180) scale(-1, 1)" d="" />
<path transform="translate(1916.000000) rotate(180) scale(-1, 1)" d="M745 538h-589v151h589v-151z" />
<path transform="translate(2802.000000) rotate(180) scale(-1, 1)" d="M745 538h-589v151h589v-151z" />
<path transform="translate(3688.000000) rotate(180) scale(-1, 1)" d="" />
<path transform="translate(4294.000000) rotate(180) scale(-1, 1)" d="M1017 0h-817v1549h200v-1372h617v-177z" />
<path transform="translate(5345.000000) rotate(180) scale(-1, 1)" d="M151 1496q0 52 35.5 87.5t87.5 35.5q53 0 89.5 -35t36.5 -88q0 -52 -36.5 -86.5t-89.5 -34.5t-88 34.5t-35 86.5zM174 0v1106h196v-1106h-196z" />
<path transform="translate(5890.000000) rotate(180) scale(-1, 1)" d="M1110 0h-195v629q0 341 -249 341q-129 0 -212 -96t-83 -243v-631h-196v1106h196v-183h4q125 209 363 209q182 0 277 -116.5t95 -339.5v-676z" />
<path transform="translate(7152.000000) rotate(180) scale(-1, 1)" d="M1115 0h-267l-473 533h-4v-533h-196v1637h196v-1038h4l449 507h252l-500 -532z" />
</g>`)
	}
}
