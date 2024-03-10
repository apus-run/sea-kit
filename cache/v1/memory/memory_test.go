package memory

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	v1 "github.com/apus-run/sea-kit/cache/v1"
)

func TestMemory_All(t *testing.T) {
	Convey("test get client", t, func() {
		mc := NewCache()
		ctx := context.Background()

		Convey("string get set", func() {
			err := mc.Set(ctx, "foo", "bar", 1*time.Hour)
			So(err, ShouldBeNil)
			val, err := mc.Get(ctx, "foo")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, "bar")

			err = mc.SetTTL(ctx, "foo", 1*time.Minute)
			So(err, ShouldBeNil)
			du, err := mc.GetTTL(ctx, "foo")
			So(err, ShouldBeNil)
			So(du, ShouldBeLessThanOrEqualTo, 1*time.Minute)
			err = mc.Del(ctx, "foo")
			So(err, ShouldBeNil)
			val, err = mc.Get(ctx, "foo")
			So(err, ShouldEqual, v1.ErrKeyNotFound)

		})

		Convey("obj get set", func() {
			type Bar struct {
				Name string
			}
			obj := Bar{
				Name: "bar",
			}
			err := mc.SetObj(ctx, "foo", obj, 1*time.Hour)
			So(err, ShouldBeNil)
			objNew := Bar{}
			err = mc.GetObj(ctx, "foo", &objNew)
			So(err, ShouldBeNil)
			So(objNew.Name, ShouldEqual, "bar")
			err = mc.Del(ctx, "foo")
			So(err, ShouldBeNil)
		})

		Convey("many op", func() {
			err := mc.SetMany(ctx, map[string]string{
				"foo1": "bar1",
				"foo2": "bar2",
			}, 1*time.Hour)
			So(err, ShouldBeNil)

			ret, err := mc.GetMany(ctx, []string{"foo1", "foo2"})
			So(err, ShouldBeNil)
			So(len(ret), ShouldEqual, 2)
			So(ret, ShouldContainKey, "foo2")
			So(ret["foo2"], ShouldEqual, "bar2")

			err = mc.DelMany(ctx, []string{"foo1", "foo2"})
			So(err, ShouldBeNil)
		})

		Convey("calc op", func() {
			val, err := mc.Increment(ctx, "foo")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 1)
			val, err = mc.Calc(ctx, "foo", 2)
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 3)
			val, err = mc.Decrement(ctx, "foo")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 2)
			err = mc.Del(ctx, "foo")
			So(err, ShouldBeNil)
		})

	})
}
