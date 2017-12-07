package churnPredict

//
//#include <string.h>
//#include <stdlib.h>
//#include <math.h>
//
//#define STRING_EQUAL(a,b) (strcasecmp(a,b)==0)
//#define STRING_NOT_EQUAL(a,b) (strcasecmp(a,b)!=0)
//#define STRING_SET(a,b)   strcpy(a,b)
//
//
//void _GLZ_Prediction_T16_36_18(
//      const double * _TRIAL_1ST_7_DAYS_CONNECTIONS__,
//      const double * _TRIAL_1ST_14_DAYS_CONNECTIONS__,
//      const double * _TRIAL_1ST_21_DAYS_CONNECTIONS__,
//      const double * _TRIAL_TOTAL_CONNECTIONS__,
//      const double * _FIRST_7_DAYS_CONNECTIONS__,
//      const double * _FIRST_14_DAYS_CONNECTIONS__,
//      const double * _FIRST_21_DAYS_CONNECTIONS__,
//      const double * _FIRST_MONTH_CONNECTIONS__,
//      const double * _MAX_COMPUTER_LIMIT__,
//      const char * _LOCALE__,
//      const char * _TRIAL_TYPE__,
//      char * _pRet
//   )
//   {
//    double _betas[14];
//    const char* _Factors[2];
//    double _Covariates[9];
//    char* _DMatrix[11][14];
//
//    int _dumcode[2];
//    int _nparam;
//    int _npco;
//    int _npca;
//    double _probs[2];
//    int _ndcats;
//    int _maxl;
//    double _maxval;
//    double _val=0;
//    int _FactorID;
//    double _pval;
//    int _i=0, _j=0, _k=0, _ifind=0;
//    int _l=0;
//    const char* _DCats[] = { "Cancelled", "Renewed" };
//
//    *_pRet=0;
//
//    for(_i=0; _i <11;_i++){
//      for(_j=0; _j <14;_j++){
//        _DMatrix[_i][_j] = "";
//      }
//    }
//
//    _betas[0] = -6.49863593786570e-002;
//    _betas[1] = 5.03464472645884e-003;
//    _betas[2] = -1.45152789513434e-001;
//    _betas[3] = 4.64166126365846e-002;
//    _betas[4] = -5.45820709343422e-003;
//    _betas[5] = 3.45246538267096e-002;
//    _betas[6] = -5.39259071666586e-002;
//    _betas[7] = 3.75987060602173e-002;
//    _betas[8] = 3.82949686933930e-002;
//    _betas[9] = 3.90748374590203e-002;
//    _betas[10] = 2.66075150294087e-001;
//    _betas[11] = 5.67108521337961e-003;
//    _betas[12] = 7.00645951922418e-001;
//    _betas[13] = 1.67135459498056e-001;
//
//    _Factors[0] = _LOCALE__;
//    _Factors[1] = _TRIAL_TYPE__;
//
//    _Covariates[0] = *_TRIAL_1ST_7_DAYS_CONNECTIONS__;
//    _Covariates[1] = *_TRIAL_1ST_14_DAYS_CONNECTIONS__;
//    _Covariates[2] = *_TRIAL_1ST_21_DAYS_CONNECTIONS__;
//    _Covariates[3] = *_TRIAL_TOTAL_CONNECTIONS__;
//    _Covariates[4] = *_FIRST_7_DAYS_CONNECTIONS__;
//    _Covariates[5] = *_FIRST_14_DAYS_CONNECTIONS__;
//    _Covariates[6] = *_FIRST_21_DAYS_CONNECTIONS__;
//    _Covariates[7] = *_FIRST_MONTH_CONNECTIONS__;
//    _Covariates[8] = *_MAX_COMPUTER_LIMIT__;
//    _DMatrix[0][1]="1";
//    _DMatrix[1][2]="1";
//    _DMatrix[2][3]="1";
//    _DMatrix[3][4]="1";
//    _DMatrix[4][5]="1";
//    _DMatrix[5][6]="1";
//    _DMatrix[6][7]="1";
//    _DMatrix[7][8]="1";
//    _DMatrix[8][9]="1";
//    _DMatrix[9][10]="en_US";
//    _DMatrix[9][11]="en_GB";
//    _DMatrix[9][12]="en_CA";
//    _DMatrix[10][13]="Trial";
//
//    _dumcode[0]=1;
//    _dumcode[1]=-1;
//    _ndcats=2;
//
//    for(_i=0;_i<2;_i++)
//    {
//        _probs[_i] = 0;
//    }
//
//    _nparam = 14;
//    _npco=9;
//    _npca=2;
//    _maxl=0;
//    _maxval=-1.0e+307;
//    if( isnan(*_TRIAL_1ST_7_DAYS_CONNECTIONS__) || isnan(*_TRIAL_1ST_14_DAYS_CONNECTIONS__) || isnan(*_TRIAL_1ST_21_DAYS_CONNECTIONS__) || isnan(*_TRIAL_TOTAL_CONNECTIONS__) || isnan(*_FIRST_7_DAYS_CONNECTIONS__) || isnan(*_FIRST_14_DAYS_CONNECTIONS__) || isnan(*_FIRST_21_DAYS_CONNECTIONS__) || isnan(*_FIRST_MONTH_CONNECTIONS__) || isnan(*_MAX_COMPUTER_LIMIT__) ||  STRING_EQUAL(_LOCALE__,"")  ||  STRING_EQUAL(_TRIAL_TYPE__,"")  ){
//      *_pRet=0;
//      goto EndProgram;
//    }
//
//    for(_l=0;_l <_ndcats;_l++){
//      _val=0;
//      if(_l != _ndcats-1){
//        for(_i=0;_i <_nparam;_i++){
//          _pval=1;
//          for(_j=0;_j <_npco+_npca;_j++){
//            if(STRING_NOT_EQUAL(_DMatrix[_j][_i],"")){
//              if(_j  <_npco){
//                _pval=_pval* (pow(_Covariates[_j] , atof(_DMatrix[_j][_i])));
//              }
//              else{
//                _FactorID=_j-_npco;
//                if(STRING_EQUAL(_DMatrix[_j][_i],_Factors[_FactorID])){
//                  _pval=_pval*_dumcode[0];
//                }
//                else{
//                  if(_dumcode[1]==-1){/*sigma-restricted*/
//                    _ifind=0;
//                    for(_k=0; _k <_nparam;_k++){
//                      if(STRING_EQUAL(_DMatrix[_j][_k],_Factors[_FactorID])){
//                        _ifind=1;
//                      }
//                    }
//                    if(_ifind==1){
//                      _pval=0;
//                    }
//                    else{
//                      _pval = _pval* _dumcode[1];
//                    }
//                  }
//                  else{
//                    _pval=_pval*_dumcode[1];
//                  }
//                }
//              }
//              }
//            }
//            _val = _val + _betas[_i+(_l)*_nparam]*_pval;
//        }
//      }
//      _probs[_l]=_val;
//      if(_val > _maxval){
//         _maxl=_l;
//        _maxval=_val;
//      }
//    }
//
//
//    _probs[0]=exp(_probs[0])/(1.0+exp(_probs[0]));
//    _probs[1]=1.0-_probs[0];
//
//
//    STRING_SET(_pRet,_DCats[_maxl]);
//
//    EndProgram:;
//
//}
//
import "C"

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log Is the Default package logger
var log = logger.GetLogger("activity-tibco-inference")

// InferfenceActivity is an Activity that is used to Invoke a ML Model using flogo-ml framework
type ModelActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a New InferfenceActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &ModelActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *ModelActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Runs an ML model
func (a *ModelActivity) Eval(context activity.Context) (done bool, err error) {

	var i1 = C.double(float64(context.GetInput("_TRIAL_1ST_7_DAYS_CONNECTIONS__").(float64)))
	var i2 = C.double(float64(context.GetInput("_TRIAL_1ST_14_DAYS_CONNECTIONS__").(float64)))
	var i3 = C.double(float64(context.GetInput("_TRIAL_1ST_21_DAYS_CONNECTIONS__").(float64)))
	var i4 = C.double(float64(context.GetInput("_TRIAL_TOTAL_CONNECTIONS__").(float64)))
	var i5 = C.double(float64(context.GetInput("_FIRST_7_DAYS_CONNECTIONS__").(float64)))
	var i6 = C.double(float64(context.GetInput("_FIRST_14_DAYS_CONNECTIONS__").(float64)))
	var i7 = C.double(float64(context.GetInput("_FIRST_21_DAYS_CONNECTIONS__").(float64)))
	var i8 = C.double(float64(context.GetInput("_FIRST_MONTH_CONNECTIONS__").(float64)))
	var i9 = C.double(float64(context.GetInput("_MAX_COMPUTER_LIMIT__").(float64)))
	var i10 = C.CString(context.GetInput("_LOCALE__").(string))
	var i11 = C.CString(context.GetInput("_TRIAL_TYPE__").(string))

	var result C.char
	C._GLZ_Prediction_T16_36_18(&i1, &i2, &i3, &i4, &i5, &i6, &i7, &i8, &i9, i10, i11, &result)

	context.SetOutput("result", C.GoString(&result))

	return true, nil
}
