package uikit

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestModernWidgets(t *testing.T) {
	test.NewApp()

	// 1. ModernTextBox
	tb := NewTextBox("Título", "placeholder", "Paisaje", true, nil)
	if tb.Text() != "Paisaje" {
		t.Errorf("Expected 'Paisaje', got '%s'", tb.Text())
	}
	tb.SetFloat(1250.50, 2)
	if tb.Float() != 1250.50 {
		t.Errorf("Expected 1250.50, got %f", tb.Float())
	}
	tb.SetInt(1995)
	if tb.Int() != 1995 {
		t.Errorf("Expected 1995, got %d", tb.Int())
	}

	// CanvasObject implementation
	var obj fyne.CanvasObject = tb
	if obj.MinSize().Height <= 0 {
		t.Errorf("Expected positive MinSize.Height")
	}

	// 2. FieldGroup batch toggle
	tb2 := NewUpperTextBox("ID", "ej: 1-01", "", true, nil)
	group := FieldGroup{tb, tb2}
	group.SetEnabled(true)
	if tb.Entry.Disabled() || tb2.Entry.Disabled() {
		t.Errorf("Expected both to be enabled")
	}
	group.SetEnabled(false)
	if !tb.Entry.Disabled() || !tb2.Entry.Disabled() {
		t.Errorf("Expected both to be disabled")
	}

	// 3. ModernBadge
	badge := NewBadge("Success", "Vendible", false)
	if badge.Text() != "Vendible" {
		t.Errorf("Expected Vendible, got %s", badge.Text())
	}

	// 3.1 Soporte widget.Importance estándar de Fyne
	badgeImp := NewBadge(widget.DangerImportance, "Error Crítico")
	if badgeImp.Text() != "Error Crítico" || badgeImp.variant != BadgeDanger {
		t.Errorf("Expected Danger variant from widget.Importance, got variant %d", badgeImp.variant)
	}
	badgeImp2 := NewBadge("Completado", widget.SuccessImportance)
	if badgeImp2.variant != BadgeSuccess {
		t.Errorf("Expected Success variant from widget.Importance, got variant %d", badgeImp2.variant)
	}

	// 3.2 Soporte de Alignment horizontal
	badge.SetAlignment(fyne.TextAlignTrailing)
	if align, ok := badge.Alignment(); !ok || align != fyne.TextAlignTrailing {
		t.Errorf("Expected alignment Trailing, got ok=%v, align=%v", ok, align)
	}

	// 4. ModernCard
	card := NewCard(tb)
	if card.cornerRadius != 10 {
		t.Errorf("Expected cornerRadius 10, got %f", card.cornerRadius)
	}
	card.WithRadius(14).WithBorder("#bfdbfe", 2).WithBackground("#ffffff")
	if card.cornerRadius != 14 || card.strokeWidth != 2 {
		t.Errorf("Expected updated radius and stroke width, got r=%f, w=%f", card.cornerRadius, card.strokeWidth)
	}
	var cardObj fyne.CanvasObject = card
	if cardObj.MinSize().Height <= 0 {
		t.Errorf("Expected positive card MinSize.Height")
	}

	// 4.1 Comprobar sombras Fyne 2.8 en ModernCard
	if !card.hasShadow {
		t.Errorf("Expected card to have shadow by default with elevation 1")
	}
	if card.bg.Shadow.BlurRadius != 6 {
		t.Errorf("Expected card shadow BlurRadius 6, got %f", card.bg.Shadow.BlurRadius)
	}
	card.WithElevation(2)
	if card.bg.Shadow.BlurRadius != 10 {
		t.Errorf("Expected card shadow BlurRadius 10 at elevation 2, got %f", card.bg.Shadow.BlurRadius)
	}
	card.WithoutShadow()
	if card.hasShadow || card.bg.Shadow.BlurRadius != 0 {
		t.Errorf("Expected shadow to be disabled with WithoutShadow()")
	}
}

func TestModernCard_PaddingConfiguration(t *testing.T) {
	app := test.NewApp()
	th := NewAppTheme()
	app.Settings().SetTheme(th)

	tb := NewTextBox("Campo", "holder", "Valor", false, nil)

	// 1. Por defecto en NewAppTheme(), una tarjeta estándar toma el padding configurado (18, 20, 18, 20)
	cardDefault := NewCard(tb)
	top, right, bottom, left := cardDefault.effectivePadding()
	if top != 18 || bottom != 18 || left != 20 || right != 20 {
		t.Errorf("Expected default padding (18, 20, 18, 20), got (%f, %f, %f, %f)", top, right, bottom, left)
	}

	// 2. Personalización explícita en la tarjeta
	cardCustom := NewTintedCard(tb, "#0f172a", "#1e293b").WithVerticalPadding(22, 10)
	top, right, bottom, left = cardCustom.effectivePadding()
	if top != 22 || bottom != 10 {
		t.Errorf("Expected vertical padding (22, 10), got top=%f, bottom=%f", top, bottom)
	}
	if left != 20 || right != 20 {
		t.Errorf("Expected lateral padding to remain 20, got left=%f, right=%f", left, right)
	}

	// 3. Modificación del padding por defecto en el tema
	th.SetCardVerticalPadding(24, 12)
	cardThemed := NewCard(tb)
	top, right, bottom, left = cardThemed.effectivePadding()
	if top != 24 || bottom != 12 {
		t.Errorf("Expected theme vertical padding (24, 12), got top=%f, bottom=%f", top, bottom)
	}

	// 4. La configuración explícita de la tarjeta tiene prioridad sobre el tema
	cardOverride := NewCard(tb).WithPaddingTop(30)
	top, right, bottom, left = cardOverride.effectivePadding()
	if top != 30 {
		t.Errorf("Expected card override top=30, got %f", top)
	}
	if bottom != 12 {
		t.Errorf("Expected theme bottom=12, got %f", bottom)
	}
}
