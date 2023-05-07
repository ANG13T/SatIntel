package osint import (
	"math"
)

// this procedure initializes variables for sgp4.
func sgp4init(opsmode *string, epoch float64, satrec *Satellite) (position, velocity Vector3) {
	var cc1sq, cc2, cc3, coef, coef1, cosio4, eeta, etasq, perige, pinvsq, psisq, qzms24, sfour, temp, temp1, temp2, temp3, temp4, tsi, xhdot1 float64

	// Deep space vars
	var cosim, sinim, em, emsq, argpm, nodem, inclm, mm, nm, s1, s2, s3, s4, s5, ss1, ss2, ss3, ss4, ss5, sz1, sz3, sz11, sz13, sz21, sz23, sz31, sz33, tc, z1, z3, z11, z13, z21, z23, z31, z33, xpidot float64

	satrec.method = "n"
	satrec.operationmode = "i"

	radiusearthkm := satrec.whichconst.radiusearthkm
	j2 := satrec.whichconst.j2
	j4 := satrec.whichconst.j4
	j3oj2 := satrec.whichconst.j3oj2

	ss := 78.0/radiusearthkm + 1.0
	qzms2ttemp := (120.0 - 78.0) / radiusearthkm
	qzms2t := qzms2ttemp * qzms2ttemp * qzms2ttemp * qzms2ttemp
	x2o3 := 2.0 / 3.0

	satrec.init = "y"
	satrec.t = 0.0

	var _, no, ao, con41, con42, cosio, cosio2, eccsq, omeosq, posq, rp, rteosq, sinio, gsto = initl(satrec.satnum, satrec.whichconst, satrec.ecco, epoch, satrec.inclo, satrec.no, satrec.method, satrec.operationmode)

	satrec.no = no
	satrec.con41 = con41
	satrec.gsto = gsto

	satrec.Error = 0

	if omeosq >= 0.0 || satrec.no >= 0.0 {
		satrec.isimp = 0
		if rp < 220.0/radiusearthkm+1.0 {
			satrec.isimp = 1
		}
		sfour = ss
		qzms24 = qzms2t
		perige = (rp - 1.0) * radiusearthkm

		if perige < 156.0 {
			sfour = perige - 78.0
			if perige < 98.0 {
				sfour = 20.0
			}

			qzms24temp := (128.0 - sfour) / radiusearthkm
			qzms24 = qzms24temp * qzms24temp * qzms24temp * qzms24temp
			sfour = sfour/radiusearthkm + 1.0
		}

		pinvsq = 1.0 / posq

		tsi = 1.0 / (ao - sfour)
		satrec.eta = ao * satrec.ecco * tsi
		etasq = satrec.eta * satrec.eta
		eeta = satrec.ecco * satrec.eta
		psisq = math.Abs(1.0 - etasq)
		coef = qzms24 * math.Pow(tsi, 4.0)
		coef1 = coef / math.Pow(psisq, 3.5)
		cc2 = coef1 * satrec.no * (ao*(1.0+1.5*etasq+eeta*(4.0+etasq)) + 0.375*j2*tsi/psisq*satrec.con41*(8.0+3.0*etasq*(8.0+etasq)))
		satrec.cc1 = satrec.bstar * cc2
		cc3 = 0.0

		if satrec.ecco > 1.0e-4 {
			cc3 = -2.0 * coef * tsi * j3oj2 * satrec.no * sinio / satrec.ecco
		}

		satrec.x1mth2 = 1.0 - cosio2
		satrec.cc4 = 2.0 * satrec.no * coef1 * ao * omeosq * (satrec.eta*(2.0+0.5*etasq) + satrec.ecco*(0.5+2.0*etasq) - j2*tsi/(ao*psisq)*(-3.0*satrec.con41*(1.0-2.0*eeta+etasq*(1.5-0.5*eeta))+0.75*satrec.x1mth2*(2.0*etasq-eeta*(1.0+etasq))*math.Cos(2.0*satrec.argpo)))
		satrec.cc5 = 2.0 * coef1 * ao * omeosq * (1.0 + 2.75*(etasq+eeta) + eeta*etasq)
		cosio4 = cosio2 * cosio2
		temp1 = 1.5 * j2 * pinvsq * satrec.no
		temp2 = 0.5 * temp1 * j2 * pinvsq
		temp3 = -0.46875 * j4 * pinvsq * pinvsq * satrec.no

		satrec.mdot = satrec.no + 0.5*temp1*rteosq*satrec.con41 + 0.0625*temp2*rteosq*(13.0-78.0*cosio2+137.0*cosio4)
		satrec.argpdot = (-0.5*temp1*con42 + 0.0625*temp2*(7.0-114.0*cosio2+395.0*cosio4) + temp3*(3.0-36.0*cosio2+49.0*cosio4))
		xhdot1 = -temp1 * cosio
		satrec.nodedot = xhdot1 + (0.5*temp2*(4.0-19.0*cosio2)+2.0*temp3*(3.0-7.0*cosio2))*cosio

		xpidot = satrec.argpdot + satrec.nodedot

		satrec.omgcof = satrec.bstar * cc3 * math.Cos(satrec.argpo)
		satrec.xmcof = 0.0

		if satrec.ecco > 1.0e-4 {
			satrec.xmcof = -x2o3 * coef * satrec.bstar / eeta
		}

		satrec.nodecf = 3.5 * omeosq * xhdot1 * satrec.cc1
		satrec.t2cof = 1.5 * satrec.cc1

		if math.Abs(cosio+1.0) > 1.5e-12 {
			satrec.xlcof = -0.25 * j3oj2 * sinio * (3.0 + 5.0*cosio) / (1.0 + cosio)
		} else {
			satrec.xlcof = -0.25 * j3oj2 * sinio * (3.0 + 5.0*cosio) / temp4
		}

		satrec.aycof = -0.5 * j3oj2 * sinio
		delmotemp := 1.0 + satrec.eta*math.Cos(satrec.mo)
		satrec.delmo = delmotemp * delmotemp * delmotemp
		satrec.sinmao = math.Sin(satrec.mo)
		satrec.x7thm1 = 7.0*cosio2 - 1.0

		if 2*math.Pi/satrec.no >= 225.0 {
			satrec.method = "d"
			satrec.isimp = 1
			tc = 0.0
			inclm = satrec.inclo

			dscomResults := dscom(epoch, satrec.ecco, satrec.argpo, tc, satrec.inclo, satrec.nodeo, satrec.no, satrec.e3, satrec.ee2, satrec.peo, satrec.pgho, satrec.pho, satrec.pinco, satrec.plo, satrec.se2, satrec.se3, satrec.sgh2, satrec.sgh3, satrec.sgh4, satrec.sh2, satrec.sh3, satrec.si2, satrec.si3, satrec.sl2, satrec.sl3, satrec.sl4, satrec.xgh2, satrec.xgh3, satrec.xgh4, satrec.xh2, satrec.xh3, satrec.xi2, satrec.xi3, satrec.xl2, satrec.xl3, satrec.xl4, satrec.zmol, satrec.zmos)

			sinim = dscomResults.sinim
			cosim = dscomResults.cosim
			satrec.e3 = dscomResults.e3
			satrec.ee2 = dscomResults.ee2
			em = dscomResults.em
			emsq = dscomResults.emsq
			satrec.peo = dscomResults.peo
			satrec.pgho = dscomResults.pgho
			satrec.pho = dscomResults.pho
			satrec.pinco = dscomResults.pinco
			satrec.plo = dscomResults.plo
			satrec.se2 = dscomResults.se2
			satrec.se3 = dscomResults.se3
			satrec.sgh2 = dscomResults.sgh2
			satrec.sgh3 = dscomResults.sgh3
			satrec.sgh4 = dscomResults.sgh4
			satrec.sh2 = dscomResults.sh2
			satrec.sh3 = dscomResults.sh3
			satrec.si2 = dscomResults.si2
			satrec.si3 = dscomResults.si3
			satrec.sl2 = dscomResults.sl2
			satrec.sl3 = dscomResults.sl3
			satrec.sl4 = dscomResults.sl4
			s1 = dscomResults.s1
			s2 = dscomResults.s2
			s3 = dscomResults.s3
			s4 = dscomResults.s4
			s5 = dscomResults.s5
			ss1 = dscomResults.ss1
			ss2 = dscomResults.ss2
			ss3 = dscomResults.ss3
			ss4 = dscomResults.ss4
			ss5 = dscomResults.ss5
			sz1 = dscomResults.sz1
			sz3 = dscomResults.sz3
			sz11 = dscomResults.sz11
			sz13 = dscomResults.sz13
			sz21 = dscomResults.sz21
			sz23 = dscomResults.sz23
			sz31 = dscomResults.sz31
			sz33 = dscomResults.sz33
			satrec.xgh2 = dscomResults.xgh2
			satrec.xgh3 = dscomResults.xgh3
			satrec.xgh4 = dscomResults.xgh4
			satrec.xh2 = dscomResults.xh2
			satrec.xh3 = dscomResults.xh3
			satrec.xi2 = dscomResults.xi2
			satrec.xi3 = dscomResults.xi3
			satrec.xl2 = dscomResults.xl2
			satrec.xl3 = dscomResults.xl3
			satrec.xl4 = dscomResults.xl4
			nm = dscomResults.nm
			z1 = dscomResults.z1
			z3 = dscomResults.z3
			z11 = dscomResults.z11
			z13 = dscomResults.z13
			z21 = dscomResults.z21
			z23 = dscomResults.z23
			z31 = dscomResults.z31
			z33 = dscomResults.z33
			satrec.zmol = dscomResults.zmol
			satrec.zmos = dscomResults.zmos

			dpperResults := dpper(satrec, inclm, satrec.init, satrec.ecco, satrec.inclo, satrec.nodeo, satrec.argpo, satrec.mo, satrec.operationmode)

			satrec.ecco = dpperResults.ep
			satrec.inclo = dpperResults.inclp
			satrec.nodeo = dpperResults.nodep
			satrec.argpo = dpperResults.argpp
			satrec.mo = dpperResults.mp

			argpm = 0.0
			nodem = 0.0
			mm = 0.0

			dsinitResults := dsinit(satrec.whichconst, cosim, emsq, satrec.argpo, s1, s2, s3, s4, s5, sinim, ss1, ss2, ss3, ss4, ss5, sz1, sz3, sz11, sz13, sz21, sz23, sz31, sz33, satrec.t, tc, satrec.gsto, satrec.mo, satrec.mdot, satrec.no, satrec.nodeo, satrec.nodedot, xpidot, z1, z3, z11, z13, z21, z23, z31, z33, satrec.ecco, eccsq, em, argpm, inclm, mm, nm, nodem, satrec.irez, satrec.atime, satrec.d2201, satrec.d2211, satrec.d3210, satrec.d3222, satrec.d4410, satrec.d4422, satrec.d5220, satrec.d5232, satrec.d5421, satrec.d5433, satrec.dedt, satrec.didt, satrec.dmdt, satrec.dnodt, satrec.domdt, satrec.del1, satrec.del2, satrec.del3, satrec.xfact, satrec.xlamo, satrec.xli, satrec.xni)

			em = dsinitResults.em
			argpm = dsinitResults.argpm
			inclm = dsinitResults.inclm
			mm = dsinitResults.mm
			nm = dsinitResults.nm
			nodem = dsinitResults.nodem
			satrec.irez = dsinitResults.irez
			satrec.atime = dsinitResults.atime
			satrec.d2201 = dsinitResults.d2201
			satrec.d2211 = dsinitResults.d2211
			satrec.d3210 = dsinitResults.d3210
			satrec.d3222 = dsinitResults.d3222
			satrec.d4410 = dsinitResults.d4410
			satrec.d4422 = dsinitResults.d4422
			satrec.d5220 = dsinitResults.d5220
			satrec.d5232 = dsinitResults.d5232
			satrec.d5421 = dsinitResults.d5421
			satrec.d5433 = dsinitResults.d5433
			satrec.dedt = dsinitResults.dedt
			satrec.didt = dsinitResults.didt
			satrec.dmdt = dsinitResults.dmdt
			satrec.dnodt = dsinitResults.dnodt
			satrec.domdt = dsinitResults.domdt
			satrec.del1 = dsinitResults.del1
			satrec.del2 = dsinitResults.del2
			satrec.del3 = dsinitResults.del3
			satrec.xfact = dsinitResults.xfact
			satrec.xlamo = dsinitResults.xlamo
			satrec.xli = dsinitResults.xli
			satrec.xni = dsinitResults.xni
		}

		if satrec.isimp != 1 {
			cc1sq = satrec.cc1 * satrec.cc1
			satrec.d2 = 4.0 * ao * tsi * cc1sq
			temp = satrec.d2 * tsi * satrec.cc1 / 3.0
			satrec.d3 = (17.0*ao + sfour) * temp
			satrec.d4 = 0.5 * temp * ao * tsi * (221.0*ao + 31.0*sfour) * satrec.cc1
			satrec.t3cof = satrec.d2 + 2.0*cc1sq
			satrec.t4cof = 0.25 * (3.0*satrec.d3 + satrec.cc1*(12.0*satrec.d2+10.0*cc1sq))
			satrec.t5cof = 0.2 * (3.0*satrec.d4 + 12.0*satrec.cc1*satrec.d3 + 6.0*satrec.d2*satrec.d2 + 15.0*cc1sq*(2.0*satrec.d2+cc1sq))
		}
	}

	position, velocity = sgp4(satrec, 0.0)
	satrec.init = "n"

	return
}

// this procedure initializes the spg4 propagator. all the initialization is consolidated here instead of having multiple loops inside other routines.
func initl(satn int64, grav GravConst, ecco, epoch, inclo, noIn float64, methodIn, opsmode string) (ainv, no, ao, con41, con42, cosio, cosio2, eccsq, omeosq, posq, rp, rteosq, sinio, gsto float64) {
	var ak, d1, adel, po float64

	x2o3 := 2.0 / 3.0

	eccsq = ecco * ecco
	omeosq = 1.0 - eccsq
	rteosq = math.Sqrt(omeosq)
	cosio = math.Cos(inclo)
	cosio2 = cosio * cosio

	ak = math.Pow(grav.xke/noIn, x2o3)
	d1 = 0.75 * grav.j2 * (3.0*cosio2 - 1.0) / (rteosq * omeosq)
	del_ := d1 / (ak * ak)
	adel = ak * (1.0 - del_*del_ - del_*(1.0/3.0+134.0*del_*del_/81.0))
	del_ = d1 / (adel * adel)
	no = noIn / (1.0 + del_)

	ao = math.Pow(grav.xke/no, x2o3)
	sinio = math.Sin(inclo)
	po = ao * omeosq
	con42 = 1.0 - 5.0*cosio2
	con41 = -con42 - cosio2 - cosio2
	ainv = 1.0 / ao
	posq = po * po
	rp = ao * (1.0 - ecco)

	if opsmode == "a" {
		ts70 := epoch - 7305.0
		ds70 := math.Floor(ts70 - 1.0e-8)
		tfrac := ts70 - ds70
		c1 := 1.72027916940703639e-2
		thgr70 := 1.7321343856509374
		fk5r := 5.07551419432269442e-15
		c1p2p := c1 + TWOPI
		gsto = math.Mod((thgr70 + c1*ds70 + c1p2p*tfrac + ts70*ts70*fk5r), TWOPI)
		if gsto < 0.0 {
			gsto = gsto + TWOPI
		}
	} else {
		gsto = gstime(epoch + 2433281.5)
	}

	return
}

// Calculates position and velocity vectors for given time
func Propagate(sat Satellite, year int, month int, day, hours, minutes, seconds int) (position, velocity Vector3) {
	j := JDay(year, month, day, hours, minutes, seconds)
	m := (j - sat.jdsatepoch) * 1440
	return sgp4(&sat, m)
}

// this procedure is the sgp4 prediction model from space command. this is an updated and combined version of sgp4 and sdp4, which were originally published separately in spacetrack report #3. this version follows the methodology from the aiaa paper (2006) describing the history and development of the code.
// satrec - initialized Satellite struct from sgp4init
// tsince - time since epoch in minutes
func sgp4(satrec *Satellite, tsince float64) (position, velocity Vector3) {
	var am, axnl, aynl, betal, cosim, sinim, cnod, snod, cos2u, sin2u, coseo1, sineo1, cosi, sini, cosip, sinip, cosisq, cossu, sinsu, cosu, sinu, delm, delomg, emsq, ecose, el2, eo1, esine, argpm, argpp, pl, rdotl, rl, rvdot, rvdotl, su, t2, t3, t4, tc, tem5, temp, temp1, temp2, tempa, tempe, templ, u, ux, uy, uz, vx, vy, vz, inclm, mm, nm, nodem, xinc, xincp, xl, xlm, mp, xmdf, xmx, xmy, nodedf, xnode, nodep, mrt float64

	mrt = 0.0
	temp4 := 1.5e-12
	x2o3 := 2.0 / 3.0

	radiusearthkm := satrec.whichconst.radiusearthkm
	xke := satrec.whichconst.xke
	j2 := satrec.whichconst.j2
	j3oj2 := satrec.whichconst.j3oj2

	vkmpersec := radiusearthkm * xke / 60.0

	satrec.t = tsince
	satrec.Error = 0
	// TODO: satrec.Error_message    = nil

	xmdf = satrec.mo + satrec.mdot*satrec.t
	var argpdf = satrec.argpo + satrec.argpdot*satrec.t
	nodedf = satrec.nodeo + satrec.nodedot*satrec.t
	argpm = argpdf
	mm = xmdf
	t2 = satrec.t * satrec.t
	nodem = nodedf + satrec.nodecf*t2
	tempa = 1.0 - satrec.cc1*satrec.t
	tempe = satrec.bstar * satrec.cc4 * satrec.t
	templ = satrec.t2cof * t2

	if satrec.isimp != 1 {
		delomg = satrec.omgcof * satrec.t
		delmtemp := 1.0 + satrec.eta*math.Cos(xmdf)
		delm = satrec.xmcof * (delmtemp*delmtemp*delmtemp - satrec.delmo)
		temp = delomg + delm
		mm = xmdf + temp
		argpm = argpdf - temp
		t3 = t2 * satrec.t
		t4 = t3 * satrec.t
		tempa = tempa - satrec.d2*t2 - satrec.d3*t3 - satrec.d4*t4
		tempe = tempe + satrec.bstar*satrec.cc5*(math.Sin(mm)-satrec.sinmao)
		templ = templ + satrec.t3cof*t3 + t4*(satrec.t4cof+satrec.t*satrec.t5cof)
	}

	nm = satrec.no
	em := satrec.ecco
	inclm = satrec.inclo

	if satrec.method == "d" {
		tc = satrec.t

		dspaceResult := dspace(satrec.irez, satrec.d2201, satrec.d2211, satrec.d3210, satrec.d3222, satrec.d4410, satrec.d4422, satrec.d5220, satrec.d5232, satrec.d5421, satrec.d5433, satrec.dedt, satrec.del1, satrec.del2, satrec.del3, satrec.didt, satrec.dmdt, satrec.dnodt, satrec.domdt, satrec.argpo, satrec.argpdot, satrec.t, tc, satrec.gsto, satrec.xfact, satrec.xlamo, satrec.no, satrec.atime, em, argpm, inclm, satrec.xli, mm, satrec.xni, nodem, nm)

		em = dspaceResult.em
		argpm = dspaceResult.argpm
		inclm = dspaceResult.inclm
		mm = dspaceResult.mm
		nodem = dspaceResult.nodem
		nm = dspaceResult.nm
	}

	if nm < 0.0 {
		satrec.Error = 2
		satrec.ErrorStr = ("Mean motion is less than zero")
	}

	am = math.Pow((xke/nm), x2o3) * tempa * tempa
	nm = xke / math.Pow(am, 1.5)
	em = em - tempe

	if em >= 1.0 || em < -0.001 {
		satrec.Error = 1
		satrec.ErrorStr = ("mean eccentricity not within range 0.0 <= e < 1.0")
	}

	if em < 1.0e-6 {
		em = 1.0e-6
	}
	mm = mm + satrec.no*templ
	xlm = mm + argpm + nodem
	emsq = em * em
	temp = 1.0 - emsq

	nodem = math.Mod(nodem, TWOPI)
	argpm = math.Mod(argpm, TWOPI)
	xlm = math.Mod(xlm, TWOPI)
	mm = math.Mod((xlm - argpm - nodem), TWOPI)

	sinim = math.Sin(inclm)
	cosim = math.Cos(inclm)

	ep := em
	xincp = inclm
	argpp = argpm
	nodep = nodem
	mp = mm
	sinip = sinim
	cosip = cosim

	if satrec.method == "d" {
		dpperResults := dpper(satrec, satrec.inclo, "n", ep, xincp, nodep, argpp, mp, satrec.operationmode)

		ep = dpperResults.ep
		xincp = dpperResults.inclp
		nodep = dpperResults.nodep
		argpp = dpperResults.argpp
		mp = dpperResults.mp

		if xincp < 0.0 {
			xincp = -xincp
			nodep = nodep + math.Pi
			argpp = argpp - math.Pi
		}

		if ep < 0.0 || ep > 1.0 {
			satrec.Error = 3
			satrec.ErrorStr = ("perturbed eccentricity not within range 0.0 <= e <= 1.0")
		}
	}

	if satrec.method == "d" {
		sinip = math.Sin(xincp)
		cosip = math.Cos(xincp)
		satrec.aycof = -0.5 * j3oj2 * sinip
		if math.Abs(cosip+1.0) > 1.5e-12 {
			satrec.xlcof = -0.25 * j3oj2 * sinip * (3.0 + 5.0*cosip) / (1.0 + cosip)
		} else {
			satrec.xlcof = -0.25 * j3oj2 * sinip * (3.0 + 5.0*cosip) / temp4
		}
	}

	axnl = ep * math.Cos(argpp)
	temp = 1.0 / (am * (1.0 - ep*ep))
	aynl = ep*math.Sin(argpp) + temp*satrec.aycof
	xl = mp + argpp + nodep + temp*satrec.xlcof*axnl

	u = math.Mod((xl - nodep), TWOPI)
	eo1 = u
	tem5 = 9999.9
	ktr := 1

	for math.Abs(tem5) >= 1.0e-12 && ktr <= 10 {
		sineo1 = math.Sin(eo1)
		coseo1 = math.Cos(eo1)
		tem5 = 1.0 - coseo1*axnl - sineo1*aynl
		tem5 = (u - aynl*coseo1 + axnl*sineo1 - eo1) / tem5
		if math.Abs(tem5) >= 0.95 {
			if tem5 > 0.0 {
				tem5 = 0.95
			} else {
				tem5 = -0.95
			}
		}
		eo1 = eo1 + tem5
		ktr = ktr + 1
	}

	ecose = axnl*coseo1 + aynl*sineo1
	esine = axnl*sineo1 - aynl*coseo1
	el2 = axnl*axnl + aynl*aynl
	pl = am * (1.0 - el2)

	if pl < 0.0 {
		satrec.Error = 4
		satrec.ErrorStr = ("semilatus rectum is less than zero")
	} else {
		rl = am * (1.0 - ecose)
		rdotl = math.Sqrt(am) * esine / rl
		rvdotl = math.Sqrt(pl) / rl
		betal = math.Sqrt(1.0 - el2)
		temp = esine / (1.0 + betal)
		sinu = am / rl * (sineo1 - aynl - axnl*temp)
		cosu = am / rl * (coseo1 - axnl + aynl*temp)
		su = math.Atan2(sinu, cosu)
		sin2u = (cosu + cosu) * sinu
		cos2u = 1.0 - 2.0*sinu*sinu
		temp = 1.0 / pl
		temp1 = 0.5 * j2 * temp
		temp2 = temp1 * temp

		if satrec.method == "d" {
			cosisq = cosip * cosip
			satrec.con41 = 3.0*cosisq - 1.0
			satrec.x1mth2 = 1.0 - cosisq
			satrec.x7thm1 = 7.0*cosisq - 1.0
		}

		mrt = rl*(1.0-1.5*temp2*betal*satrec.con41) + 0.5*temp1*satrec.x1mth2*cos2u
		su = su - 0.25*temp2*satrec.x7thm1*sin2u
		xnode = nodep + 1.5*temp2*cosip*sin2u
		xinc = xincp + 1.5*temp2*cosip*sinip*cos2u
		mvt := rdotl - nm*temp1*satrec.x1mth2*sin2u/xke
		rvdot = rvdotl + nm*temp1*(satrec.x1mth2*cos2u+1.5*satrec.con41)/xke

		sinsu = math.Sin(su)
		cossu = math.Cos(su)
		snod = math.Sin(xnode)
		cnod = math.Cos(xnode)
		sini = math.Sin(xinc)
		cosi = math.Cos(xinc)
		xmx = -snod * cosi
		xmy = cnod * cosi
		ux = xmx*sinsu + cnod*cossu
		uy = xmy*sinsu + snod*cossu
		uz = sini * sinsu
		vx = xmx*cossu - cnod*sinsu
		vy = xmy*cossu - snod*sinsu
		vz = sini * cossu

		_mr := mrt * radiusearthkm

		position.X = _mr * ux
		position.Y = _mr * uy
		position.Z = _mr * uz

		velocity.X = (mvt*ux + rvdot*vx) * vkmpersec
		velocity.Y = (mvt*uy + rvdot*vy) * vkmpersec
		velocity.Z = (mvt*uz + rvdot*vz) * vkmpersec
	}

	if mrt < 1.0 {
		satrec.Error = 6
		satrec.ErrorStr = ("mrt is less than 1.0 indicating the satellite has decayed")
	}

	return
}
