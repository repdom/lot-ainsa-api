package pdf

import (
	"be-lotsanmateo-api/internal/domain/model"
	"fmt"
	"log"
	"strconv"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

func NewPaymentPlanPDF() PaymentPlanPDF {
	return PaymentPlanPDF{}
}

type PaymentPlanPDF struct{}

type Attribute struct {
	title1 string
	value1 string
	title2 string
	value2 string
}

type PaymentPlan struct {
	Polygon   string
	Lote      string
	FullName  string
	Area      float64
	Address   string
	Phone     string
	StartDate string
	DUI       string
	Loan      model.ResponseLoan
}

func (p PaymentPlanPDF) GeneratePDF(payment PaymentPlan) ([]byte, error) {
	configuration := config.NewBuilder().
		WithPageSize(pagesize.Letter).
		WithBottomMargin(10).
		WithRightMargin(10).
		WithLeftMargin(10).
		WithMaxGridSize(100).
		Build()
	mrt := maroto.New(configuration)
	m := maroto.NewMetricsDecorator(mrt)
	log.Println("Generating PDF")
	m.AddRow(20,
		image.NewFromFileCol(25, "docs/assets/images/lot_name-rbg.png", props.Rect{
			Left:               -30,
			Percent:            78,
			Center:             true,
			JustReferenceWidth: true,
		}),
		col.New(50).Add(
			text.New(fmt.Sprintf("POLIGONO %s - LOTE %s", payment.Polygon, payment.Lote), props.Text{
				Size:  12,
				Style: fontstyle.Bold,
				Align: align.Right,
				Color: &props.Color{Green: 180},
			}),
			text.New(fmt.Sprintf("%fM2", payment.Area), props.Text{
				Style: fontstyle.Bold,
				Top:   7,
				Right: 15,
				Size:  12,
				Align: align.Right,
				Color: &props.Color{Green: 180},
			}),
		),
		image.NewFromFileCol(25, "docs/assets/images/name.png", props.Rect{
			Left: 15,
		}),
	)

	m.AddRow(5)

	m.AddRow(5,
		col.New(2),
		col.New(53).Add(
			text.New("LOLOTIQUE, SAN MIGUEL", props.Text{
				Align: align.Left,
				Style: fontstyle.Bold,
			}),
		),

		text.NewCol(20, "F. Otorgamiento:", props.Text{
			Align: align.Left,
		}),
		text.NewCol(20, payment.StartDate, props.Text{
			Align: align.Right,
		}),
	)

	addAttribute(m, Attribute{
		"Cliente:", payment.FullName, "Precio del Lote :", formatMoney(payment.Loan.Amount, ",", "."),
	}, props.Text{
		Align: align.Left,
	}, props.Text{
		Align: align.Right,
	}, &props.Cell{}, &props.Cell{
		BackgroundColor: &props.Color{Red: 182, Green: 199, Blue: 229},
	})

	addAttribute(m, Attribute{
		"Dirección:", payment.Address, "", formatMoney(0, ",", "."),
	}, props.Text{
		Align: align.Left,
	}, props.Text{
		Align: align.Right,
	}, &props.Cell{}, &props.Cell{})

	addAttribute(m, Attribute{
		"Teléfono:", payment.Phone, fmt.Sprintf("PRIMA DEL %.2f%%", payment.Loan.DownPaymentRate), formatMoney(payment.Loan.Premium, ",", "."),
	}, props.Text{
		Align: align.Left,
	}, props.Text{
		Align: align.Right,
	}, &props.Cell{}, &props.Cell{
		BackgroundColor: &props.Color{Red: 182, Green: 199, Blue: 229},
	})

	addAttribute(m, Attribute{"Lote:", payment.Lote, "", ""}, props.Text{}, props.Text{}, &props.Cell{}, &props.Cell{})

	addAttribute(m, Attribute{
		"Polígono:", payment.Polygon, fmt.Sprintf("Crédito %.2f%% :", 100-payment.Loan.DownPaymentRate), formatMoney(payment.Loan.TotalAmount, ",", "."),
	}, props.Text{
		Align: align.Left,
		Color: &props.Color{Red: 255, Green: 255, Blue: 255},
	}, props.Text{
		Align: align.Right,
		Color: &props.Color{Red: 255, Green: 255, Blue: 255},
	}, &props.Cell{
		BackgroundColor: &props.Color{Red: 255, Green: 0, Blue: 0},
	}, &props.Cell{
		BackgroundColor: &props.Color{Red: 255, Green: 0, Blue: 0},
	})

	addAttribute(m, Attribute{
		"CED", payment.DUI, "Plazo en años:", strconv.FormatFloat(payment.Loan.Years, 'f', 0, 64),
	}, props.Text{
		Align: align.Left,
	}, props.Text{
		Align: align.Right,
	}, &props.Cell{}, &props.Cell{
		BackgroundColor: &props.Color{Red: 182, Green: 199, Blue: 229},
	})

	addAttribute(m, Attribute{
		"", "", "Tasa Preferencial:", fmt.Sprintf(" %.2f%%", payment.Loan.Rate),
	}, props.Text{
		Align: align.Left,
	}, props.Text{
		Align: align.Right,
	}, &props.Cell{}, &props.Cell{})

	m.AddRow(2)
	m.AddRow(5,
		col.New(100).Add(
			text.New("Plan de pago", props.Text{
				Align: align.Center,
			}),
		),
	)

	m.AddRow(5)

	m.AddRow(5,
		col.New(2),
		col.New(14).Add(
			text.New("Capital", props.Text{
				Align: align.Center,
				Style: fontstyle.Bold,
			}),
		).
			WithStyle(&props.Cell{
				BackgroundColor: &props.Color{Red: 169, Green: 208, Blue: 142},
				BorderColor:     &props.Color{Red: 0, Green: 0, Blue: 0},
				BorderType:      border.Full,
				BorderThickness: 0.5,
			}),
		col.New(16).Add(
			text.New(formatMoney(payment.Loan.TotalAmount, ",", "."), props.Text{
				Align: align.Center,
				Color: &props.Color{Red: 255, Green: 255, Blue: 255},
			}),
		).
			WithStyle(&props.Cell{
				BackgroundColor: &props.Color{Red: 255, Green: 0, Blue: 0},
				BorderColor:     &props.Color{Red: 0, Green: 0, Blue: 0},
				BorderType:      border.Full,
				BorderThickness: 0.5,
			}),
		col.New(16).Add(
			text.New("Interés", props.Text{
				Align: align.Center,
				Style: fontstyle.Bold,
			}),
		).
			WithStyle(&props.Cell{
				BackgroundColor: &props.Color{Red: 169, Green: 208, Blue: 142},
				BorderColor:     &props.Color{Red: 0, Green: 0, Blue: 0},
				BorderType:      border.Full,
				BorderThickness: 0.5,
			}),
		col.New(10).Add(
			text.New(fmt.Sprintf("%.2f%%", payment.Loan.Rate), props.Text{
				Align: align.Center,
			}),
		).
			WithStyle(&props.Cell{
				BorderColor:     &props.Color{Red: 0, Green: 0, Blue: 0},
				BorderType:      border.Full,
				BorderThickness: 0.5,
			}),
		col.New(22).Add(
			text.New("Inicio de Préstamo", props.Text{
				Align: align.Center,
				Style: fontstyle.Bold,
			}),
		).
			WithStyle(&props.Cell{
				BackgroundColor: &props.Color{Red: 169, Green: 208, Blue: 142},
				BorderColor:     &props.Color{Red: 0, Green: 0, Blue: 0},
				BorderType:      border.Full,
				BorderThickness: 0.5,
			}),
		col.New(16).Add(
			text.New(payment.StartDate, props.Text{
				Align: align.Center,
			}),
		).
			WithStyle(&props.Cell{
				BorderColor:     &props.Color{Red: 0, Green: 0, Blue: 0},
				BorderType:      border.Full,
				BorderThickness: 0.5,
			}),
	)

	log.Printf("%d", payment.Loan.NumberOfInstallments)

	m.AddRow(5,
		col.New(2),
		col.New(14).Add(
			text.New("Años", props.Text{
				Align: align.Center,
				Style: fontstyle.Bold,
			}),
		).
			WithStyle(&props.Cell{
				BackgroundColor: &props.Color{Red: 169, Green: 208, Blue: 142},
				BorderColor:     &props.Color{Red: 0, Green: 0, Blue: 0},
				BorderType:      border.Full,
				BorderThickness: 0.5,
			}),
		col.New(16).Add(
			text.New(strconv.FormatFloat(payment.Loan.Years, 'f', 0, 64), props.Text{
				Align: align.Center,
			}),
		).
			WithStyle(&props.Cell{
				BorderColor:     &props.Color{Red: 0, Green: 0, Blue: 0},
				BorderType:      border.Full,
				BorderThickness: 0.5,
			}),
		col.New(16).Add(
			text.New("Nº de Cuotas", props.Text{
				Align: align.Center,
				Style: fontstyle.Bold,
			}),
		).
			WithStyle(&props.Cell{
				BackgroundColor: &props.Color{Red: 169, Green: 208, Blue: 142},
				BorderColor:     &props.Color{Red: 0, Green: 0, Blue: 0},
				BorderType:      border.Full,
				BorderThickness: 0.5,
			}),

		col.New(10).Add(
			text.New(strconv.Itoa(payment.Loan.NumberOfInstallments), props.Text{
				Align: align.Center,
			}),
		).
			WithStyle(&props.Cell{
				BorderColor:     &props.Color{Red: 0, Green: 0, Blue: 0},
				BorderType:      border.Full,
				BorderThickness: 0.5,
			}),
		col.New(22).Add(
			text.New("Cuota", props.Text{
				Align: align.Center,
				Style: fontstyle.Bold,
			}),
		).
			WithStyle(&props.Cell{
				BackgroundColor: &props.Color{Red: 169, Green: 208, Blue: 142},
				BorderColor:     &props.Color{Red: 0, Green: 0, Blue: 0},
				BorderType:      border.Full,
				BorderThickness: 0.5,
			}),
		col.New(16).Add(
			text.New(formatMoney(payment.Loan.MonthlyPayment, ",", "."), props.Text{
				Align: align.Center,
			}),
		).
			WithStyle(&props.Cell{
				BackgroundColor: &props.Color{Red: 255, Green: 192, Blue: 0},
				BorderColor:     &props.Color{Red: 0, Green: 0, Blue: 0},
				BorderType:      border.Full,
				BorderThickness: 0.5,
			}),
	)

	m.AddRow(5)

	err := m.RegisterHeader(
		row.New(5).Add(
			data(2, "", &props.Color{}),
			header(8, "N°"),
			header(15, "Fecha"),
			header(13, "Cuota Total"),
			header(13, "Capital"),
			header(13, "Interés"),
			header(16, "Saldo inicial"),
			header(16, "Saldo final"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	mov := payment.Loan.FeeSimulation

	for index, element := range mov {
		m.AddRow(5,
			data(2, "", &props.Color{}),
			data(8, strconv.Itoa(index+1), &props.Color{}),
			data(15, element.Payday, &props.Color{}),
			data(13, formatMoney(element.Amount, ",", "."), &props.Color{}),
			data(13, formatMoney(element.Capital, ",", "."), &props.Color{}),
			data(13, formatMoney(element.Interest, ",", "."), &props.Color{Red: 255}),
			data(16, formatMoney(element.BalanceStart, ",", "."), &props.Color{}),
			data(16, formatMoney(element.BalanceLast, ",", "."), &props.Color{}),
		)
	}

	pdf, err := m.Generate()
	return pdf.GetBytes(), err
}

func header(size int, title string) core.Col {
	return col.New(size).
		Add(
			text.New(title, props.Text{
				Align: align.Center,
				Size:  10,
			}),
		).
		WithStyle(&props.Cell{
			BackgroundColor: &props.Color{Red: 169, Green: 208, Blue: 142},
			BorderColor:     &props.Color{Red: 0, Green: 0, Blue: 0},
			BorderType:      border.Full,
			BorderThickness: 0.5,
		})
}

func data(size int, title string, color *props.Color) core.Col {
	return col.New(size).
		Add(
			text.New(title, props.Text{
				Top:   2,
				Size:  9,
				Align: align.Center,
				Color: color,
			}),
		)
}

func addAttribute(m core.Maroto, attr Attribute, ps3, ps4 props.Text, style3, style4 *props.Cell) {
	m.AddRow(5,
		col.New(2),
		col.New(18).Add(
			text.New(attr.title1, props.Text{
				Align: align.Left,
				Style: fontstyle.Bold,
			}),
		),
		col.New(35).Add(
			text.New(attr.value1, props.Text{
				Align: align.Left,
			}),
		),
		col.New(20).Add(
			text.New(attr.title2, ps3),
		).WithStyle(style3),
		col.New(20).Add(
			text.New(attr.value2, ps4),
		).WithStyle(style4),
	)
}

func formatMoney(input float64, thousand, decimal string) string {
	var result string
	var isNegative bool

	value := int(input * 100)
	if value == 0 {
		return "$                  -      "
	}

	if value < 0 {
		value = value * -1
		isNegative = true
	}
	// apply the decimal separator
	result = fmt.Sprintf("%s%02d%s", decimal, value%100, result)
	value /= 100
	// for each 3 dígits put a dot "."
	for value >= 1000 {
		result = fmt.Sprintf("%s%03d%s", thousand, value%1000, result)
		value /= 1000
	}
	if isNegative {
		return fmt.Sprintf("$ -%10d%s", value, result)
	}
	return fmt.Sprintf("$ %10d%s", value, result)
}
