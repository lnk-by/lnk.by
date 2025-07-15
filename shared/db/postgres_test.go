package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkipSQLComments_empty(t *testing.T) {
	assert.Equal(t, "", skipSQLComments(""))
}

func TestSkipSQLComments_withoutComments(t *testing.T) {
	stmt := `qwe
rty`
	assert.Equal(t, stmt, skipSQLComments(stmt))
}

func TestSkipSQLComments_onlyComments(t *testing.T) {
	stmt := `--qwe
--rty`
	assert.Equal(t, "", skipSQLComments(stmt))
}

func TestSkipSQLComments_withLeadingComment(t *testing.T) {
	stmt := `--leading comment
qwe
rty`
	expected := `qwe
rty`
	assert.Equal(t, expected, skipSQLComments(stmt))
}

func TestSkipSQLComments_withIndentedLeadingComment(t *testing.T) {
	stmt := `    --leading comment
qwe
rty`
	expected := `qwe
rty`
	assert.Equal(t, expected, skipSQLComments(stmt))
}

func TestSkipSQLComments_withTrailingComment(t *testing.T) {
	stmt := `qwe
rty
--trailing comment`
	expected := `qwe
rty`
	assert.Equal(t, expected, skipSQLComments(stmt))
}

func TestSkipSQLComments_withIndentedTrailingComment(t *testing.T) {
	stmt := `qwe
rty
    --trailing comment`
	expected := `qwe
rty`
	assert.Equal(t, expected, skipSQLComments(stmt))
}

func TestSkipSQLComments_withCommentInside(t *testing.T) {
	stmt := `qwe
--comment inside	
rty`
	expected := `qwe
rty`
	assert.Equal(t, expected, skipSQLComments(stmt))
}

func TestSkipSQLComments_withIndentedCommentInside(t *testing.T) {
	stmt := `qwe
    --comment inside
rty`
	expected := `qwe
rty`
	assert.Equal(t, expected, skipSQLComments(stmt))
}

func testSplitSQLStatements(t *testing.T, script string, expected []string) {
	actual, err := splitSQLStatements(script)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestSplitSQLStatements_empty(t *testing.T) {
	testSplitSQLStatements(t, "", nil)
}

func TestSplitSQLStatements_singleLineWithoutSemicolon(t *testing.T) {
	testSplitSQLStatements(t, "qwe", []string{"qwe"})
}

func TestSplitSQLStatements_singleLineWithColonAtEnd(t *testing.T) {
	testSplitSQLStatements(t, "qwe;", []string{"qwe"})
}

func TestSplitSQLStatements_singleLineWithSemicolonInside(t *testing.T) {
	testSplitSQLStatements(t, "qwe;rty", []string{"qwe", "rty"})
}

func TestSplitSQLStatements_multiLinesWithoutSemicolon(t *testing.T) {
	script := `qwe
rty`
	expected := []string{`qwe
rty`}
	testSplitSQLStatements(t, script, expected)
}

func TestSplitSQLStatements_multiLinesWithoutSemicolonWithLeadingSpaces(t *testing.T) {
	script := `  qwe
rty`
	expected := []string{`qwe
rty`}
	testSplitSQLStatements(t, script, expected)
}

func TestSplitSQLStatements_multiLinesWithoutSemicolonWithNewLine(t *testing.T) {
	script := `
qwe
rty`
	expected := []string{`qwe
rty`}
	testSplitSQLStatements(t, script, expected)
}

func TestSplitSQLStatements_multiLinesWithSemicolonAtEnd(t *testing.T) {
	script := `qwe
rty;`
	expected := []string{`qwe
rty`}
	testSplitSQLStatements(t, script, expected)
}

func TestSplitSQLStatements_multiLinesWithSemicolon(t *testing.T) {
	script := `qwe
rty;

asd
fgh`
	expected := []string{
		`qwe
rty`,
		`asd
fgh`,
	}
	testSplitSQLStatements(t, script, expected)
}

//

func TestSplitSQLStatements_withDollarQuote_singleLineWithoutSemicolon(t *testing.T) {
	testSplitSQLStatements(t, "qwe$zzz$rty$xxx$uio", []string{"qwe$zzz$rty$xxx$uio"})
}

func TestSplitSQLStatements_withDollarQuote_singleLineWithColonAtEnd(t *testing.T) {
	testSplitSQLStatements(t, "qwe$zzz$rty$xxx$uio;", []string{"qwe$zzz$rty$xxx$uio"})
}

func TestSplitSQLStatements_withDollarQuote_singleLineWithSemicolonInside(t *testing.T) {
	testSplitSQLStatements(t, "qwe$zzz$asd$xxx$;rty", []string{"qwe$zzz$asd$xxx$", "rty"})
}

func TestSplitSQLStatements_withDollarQuote_multiLinesWithoutSemicolon(t *testing.T) {
	script := `qwe $zzz$
rty
uio
$xxx$
asd`
	expected := []string{`qwe $zzz$
rty
uio
$xxx$
asd`}
	testSplitSQLStatements(t, script, expected)
}

func TestSplitSQLStatements_withDollarQuote_multiLinesWithSemicolon(t *testing.T) {
	script := `qwe $zzz$
rty;
uio
$xxx$
asd`
	expected := []string{`qwe $zzz$
rty;
uio
$xxx$
asd`}
	testSplitSQLStatements(t, script, expected)
}

func TestSplitSQLStatements_withDollarQuote_badScript(t *testing.T) {
	script := `qwe $zzz$
rty;
uio
asd`
	_, err := splitSQLStatements(script)
	if assert.Error(t, err) {
		assert.Equal(t, errInvalidScript, err)
	}
}
